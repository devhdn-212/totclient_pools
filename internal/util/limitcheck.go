package util

import (
	"context"
	"fmt"
	"time"

	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/redis/go-redis/v9"
)

// checkAndIncrementLimitScript atomically compares the counter at KEYS[1]
// against ARGV[2] (limit) before adding ARGV[1] (amount) — Redis runs Lua
// scripts single-threaded, so concurrent bets on the SAME key still
// serialize (same guarantee the old SELECT ... FOR UPDATE row lock gave),
// while unrelated keys never block each other.
//
// If KEYS[1] doesn't exist yet, it's seeded from ARGV[4] first (the real
// historical total from Postgres, passed in by the caller — see
// CheckAndIncrementLimitRedis) instead of assuming 0. That single missing-key
// branch covers both "genuinely never bet on" (seed is legitimately 0) and
// "Redis lost its memory of this key" (seed is whatever Postgres says was
// really wagered) — the caller can't tell those apart without asking
// Postgres, so it always passes the real seed and lets Redis decide whether
// it was actually needed.
var checkAndIncrementLimitScript = redis.NewScript(`
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
// seed is only invoked when nmCounter doesn't exist in Redis yet — see the
// Lua script comment above for why that one Postgres round-trip is
// unavoidable and why it's cheap in practice (once per key per Redis
// restart, not per bet).
func CheckAndIncrementLimitRedis(ctx context.Context, nmCounter string, amount, limit int64, ttl time.Duration, seed SeedFunc) (newTotal int64, ok bool, err error) {
	if limit <= 0 {
		return 0, true, nil
	}

	exists, err := connection.RDBLimit.Exists(ctx, nmCounter).Result()
	if err != nil {
		return 0, false, fmt.Errorf("error checking redis limit counter %s: %w", nmCounter, err)
	}

	var seedValue int64
	if exists == 0 {
		seedValue, err = seed(ctx)
		if err != nil {
			return 0, false, fmt.Errorf("error seeding redis limit counter %s from database: %w", nmCounter, err)
		}
	}

	res, err := checkAndIncrementLimitScript.Run(ctx, connection.RDBLimit,
		[]string{nmCounter}, amount, limit, int64(ttl.Seconds()), seedValue,
	).Result()
	if err != nil {
		return 0, false, fmt.Errorf("error check-and-increment redis limit counter %s: %w", nmCounter, err)
	}

	pair, valid := res.([]interface{})
	if !valid || len(pair) != 2 {
		return 0, false, fmt.Errorf("unexpected redis script result for %s: %v", nmCounter, res)
	}
	accepted, _ := pair[0].(int64)
	value, _ := pair[1].(int64)
	return value, accepted == 1, nil
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
