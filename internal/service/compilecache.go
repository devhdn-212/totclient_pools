package service

import (
	"strconv"
	"time"

	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/shopspring/decimal"
)

// compileCacheTTL is deliberately generous (a period can stay open a while,
// and the admin can take time to submit the keluaran result) but still
// bounded — once it lapses, whatever reads this cache has to fall back to
// Postgres for that period regardless (a bet still stuck "RUNNING" after
// compile runs means it was never settled, whether Redis lost it to
// expiry/eviction/restart or anything else — see DOKUMENTASI.md).
const compileCacheTTL = 2 * 24 * time.Hour

// compileDetailItem is one accepted bet line, cached so a later "compile"
// step (admin submits the keluaran result → determine win/lose for every bet
// in that period) can read it from Redis instead of querying Postgres row by
// row. Postgres remains the durable record regardless — this is written
// only after checkoutService.Submit's tx.Commit succeeds, and is pure
// acceleration, never the only copy of anything.
type compileDetailItem struct {
	Idtrxdetail   string          `json:"idtrxdetail"`
	Playerinvoice string          `json:"playerinvoice"`
	Username      string          `json:"username"`
	Nomor         string          `json:"nomor"`
	Type          string          `json:"type"`
	Posisi        string          `json:"posisi"`
	Bet           int             `json:"bet"`
	Payout        decimal.Decimal `json:"payout"`
	Win           decimal.Decimal `json:"win"`
}

func compileDetailCacheKey(idtrxkeluaran int, playerinvoice string) string {
	return "trxkeluarandetail:" + strconv.Itoa(idtrxkeluaran) + ":playerinvoice:" + playerinvoice
}

func compileListPlayerinvoiceCacheKey(idtrxkeluaran int) string {
	return "trxkeluarandetail:" + strconv.Itoa(idtrxkeluaran) + ":listplayerinvoice"
}

// cacheCompileData writes one checkout chunk's accepted bets to Redis for a
// later compile step to pick up. Meant to be called via `go` (fire-and-forget)
// right after tx.Commit succeeds — a failed cache write here is never fatal,
// compile's own RUNNING-status reconciliation against Postgres covers
// whatever Redis doesn't have.
//
// listplayerinvoice is a Redis SET (SADD, not a read-appended JSON array) —
// checkouts for the same period can commit concurrently from many
// goroutines, and SADD is atomic on its own; a read-modify-write JSON array
// would silently drop entries under that concurrency.
func cacheCompileData(idtrxkeluaran int, playerinvoice string, items []compileDetailItem) {
	if len(items) == 0 {
		return
	}
	if err := connection.SetRedis(compileDetailCacheKey(idtrxkeluaran, playerinvoice), items, compileCacheTTL); err != nil {
		return
	}
	_ = connection.AddRedisSet(compileListPlayerinvoiceCacheKey(idtrxkeluaran), playerinvoice, compileCacheTTL)
}