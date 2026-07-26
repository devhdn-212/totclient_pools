package util

import (
	"context"
	"fmt"
	"time"

	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/redis/go-redis/v9"
)

// checkAndIncrementLimitScript is the fast path — it never touches Postgres.
// It assumes KEYS[1] already exists (the common case: some bet on this
// number this period already created it) and just does the compare-then-
// increment. If the key turns out not to exist, it returns the sentinel
// {-1, 0} instead of guessing a seed, and the caller falls back to
// checkAndIncrementLimitScriptSeeded below.
//
// Redis runs Lua scripts single-threaded, so concurrent bets on the SAME
// key still serialize (same guarantee the old SELECT ... FOR UPDATE row
// lock gave), while unrelated keys never block each other.
var checkAndIncrementLimitScript = redis.NewScript(`
local exists = redis.call('EXISTS', KEYS[1])
if exists == 0 then
	return {-1, 0}
end

local current = tonumber(redis.call('GET', KEYS[1]))
local amount = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local ttlSeconds = tonumber(ARGV[3])
if ttlSeconds > 0 then
	redis.call('EXPIRE', KEYS[1], ttlSeconds)
end

if current + amount > limit then
	return {0, current}
end

local newTotal = redis.call('INCRBY', KEYS[1], amount)
return {1, newTotal}
`)

// checkAndIncrementLimitScriptSeeded is the fallback path — only run when
// the fast path above reported the key missing. KEYS[1] is seeded from
// ARGV[4] (the real historical total from Postgres, passed in by the
// caller — see CheckAndIncrementLimitRedis) instead of assuming 0, since
// "genuinely never bet on" (seed is legitimately 0) and "Redis lost its
// memory of this key" (seed is whatever Postgres says was really wagered)
// look identical from here and only Postgres can tell them apart.
//
// The EXISTS check is repeated here (not just assumed missing) because
// another request for the same key can race between the two scripts: if it
// already seeded the key in the meantime, this read the up-to-date value
// instead of re-seeding over it.
var checkAndIncrementLimitScriptSeeded = redis.NewScript(`
local exists = redis.call('EXISTS', KEYS[1])
local current
if exists == 1 then
	current = tonumber(redis.call('GET', KEYS[1]))
else
	current = tonumber(ARGV[4])
	redis.call('SET', KEYS[1], current)
end

local amount = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local ttlSeconds = tonumber(ARGV[3])
if ttlSeconds > 0 then
	redis.call('EXPIRE', KEYS[1], ttlSeconds)
end

if current + amount > limit then
	return {0, current}
end

local newTotal = redis.call('INCRBY', KEYS[1], amount)
return {1, newTotal}
`)

// SeedFunc computes the real historical total for a limit counter straight
// from the permanent bet ledger in Postgres — called only when Redis has
// never seen the key (first bet on that number this period, or Redis lost
// its memory of it, e.g. after a restart).
type SeedFunc func(ctx context.Context) (int64, error)

// limitScriptResult is the {code, value} pair every limit script returns.
// code is -1 (key missing, caller must seed and retry), 0 (rejected —
// value is the current total), or 1 (accepted — value is the new total).
func limitScriptResult(nmCounter string, res interface{}) (code, value int64, err error) {
	pair, valid := res.([]interface{})
	if !valid || len(pair) != 2 {
		return 0, 0, fmt.Errorf("unexpected redis script result for %s: %v", nmCounter, res)
	}
	code, _ = pair[0].(int64)
	value, _ = pair[1].(int64)
	return code, value, nil
}

// CheckAndIncrementLimitRedis adds amount to the running total kept at
// nmCounter in Redis and only commits the increment if the new total stays
// within limit. limit <= 0 means uncapped — always succeeds without ever
// touching Redis or calling seed (mirrors the old Postgres-backed behavior
// exactly).
//
// nmCounter already embeds the idtrxkeluaran (one per market period), so a
// new period always gets a brand new key — ttl is just a memory cleanup
// safety net for keys nobody touches again, not a correctness mechanism.
//
// Only one Redis round-trip in the common case (key already exists — true
// for any number that's had a bet on it already this period): the fast
// path above both checks and increments in that single call. seed (a
// Postgres round-trip) and the second script call only happen the first
// time a given number is ever bet on in a period, or after Redis loses its
// memory of the key (e.g. a restart) — see checkAndIncrementLimitScriptSeeded.
func CheckAndIncrementLimitRedis(ctx context.Context, nmCounter string, amount, limit int64, ttl time.Duration, seed SeedFunc) (newTotal int64, ok bool, err error) {
	if limit <= 0 {
		return 0, true, nil
	}

	ttlSeconds := int64(ttl.Seconds())

	res, err := checkAndIncrementLimitScript.Run(ctx, connection.RDBLimit,
		[]string{nmCounter}, amount, limit, ttlSeconds,
	).Result()
	if err != nil {
		return 0, false, fmt.Errorf("error check-and-increment redis limit counter %s: %w", nmCounter, err)
	}
	code, value, err := limitScriptResult(nmCounter, res)
	if err != nil {
		return 0, false, err
	}
	if code != -1 {
		return value, code == 1, nil
	}

	// Fast path reported the key missing — seed it from Postgres and retry
	// with the script that knows how to initialize it.
	seedValue, err := seed(ctx)
	if err != nil {
		return 0, false, fmt.Errorf("error seeding redis limit counter %s from database: %w", nmCounter, err)
	}

	res, err = checkAndIncrementLimitScriptSeeded.Run(ctx, connection.RDBLimit,
		[]string{nmCounter}, amount, limit, ttlSeconds, seedValue,
	).Result()
	if err != nil {
		return 0, false, fmt.Errorf("error check-and-increment redis limit counter %s (seeded): %w", nmCounter, err)
	}
	code, value, err = limitScriptResult(nmCounter, res)
	if err != nil {
		return 0, false, err
	}
	return value, code == 1, nil
}

// DecrementCounterRedis undoes a prior CheckAndIncrementLimitRedis increment
// — used when a bet passes the limittotal check but then fails limitglobal
// (limittotal's increment needs rolling back for that one bet), and when the
// surrounding checkout's Postgres transaction ends up rolling back entirely
// (every Redis increment made during that attempt needs undoing, since
// nothing about the bet actually got persisted).
func DecrementCounterRedis(ctx context.Context, nmCounter string, amount int64) error {
	if err := connection.RDBLimit.DecrBy(ctx, nmCounter, amount).Err(); err != nil {
		return fmt.Errorf("error decrement redis limit counter %s: %w", nmCounter, err)
	}
	return nil
}
