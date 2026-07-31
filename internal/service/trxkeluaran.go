package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	RedisTrxkeluaran = "agen:trxkeluaran"
	// RedisTrxkeluaranDetail mirrors the agen dashboard's own per-period
	// cache key (see totagen_api's trxkeluaran.go) — keyed by idtrxkeluaran,
	// not just idcomp. checkout.go invalidates this alongside RedisTrxkeluaran
	// since both go stale on the agen side the moment
	// total_member/total_bet/total_pairs/total_payout change here.
	RedisTrxkeluaranDetail = RedisTrxkeluaran + ":detail"

	// RedisTotalsDirty is a Redis SET of "idcomp|idtrxkeluaran" members —
	// periods whose total_member/total_bet/total_pairs/total_payout need
	// recomputing. checkout.go adds to it (MarkTotalsDirty) instead of
	// updating those columns inline, so many concurrent checkouts against
	// the same period no longer serialize on that row's lock. A background
	// ticker (see main.go) drains it periodically via FlushDirtyTotals.
	RedisTotalsDirty = "internal:totals:dirty"
)

// MarkTotalsDirty schedules idtrxkeluaran's summary columns for a background
// recompute — fire-and-forget, never blocks the checkout that calls it.
func MarkTotalsDirty(idcomp string, idtrxkeluaran int) {
	member := strings.ToLower(idcomp) + "|" + strconv.Itoa(idtrxkeluaran)
	if err := connection.RDB.SAdd(context.Background(), RedisTotalsDirty, member).Err(); err != nil {
		connection.Log.Error("MarkTotalsDirty: SAdd failed", zap.Error(err), zap.String("member", member))
	}
}

// FlushDirtyTotals drains every period marked dirty since the last call and
// recomputes its totals via TrxkeluaranRepository.RefreshTotals. Called from
// a ticker in main.go rather than inline in checkout — this is what turns N
// concurrent checkouts against one period's row lock into at most one UPDATE
// per period per tick.
//
// SPopN atomically pops (not just reads) up to count members, so a checkout
// that marks a period dirty while a flush is in flight is never lost — it
// just stays in the set for the next tick instead of being wiped by a
// read-then-delete race.
func FlushDirtyTotals(ctx context.Context, repo domain.TrxkeluaranRepository) {
	members, err := connection.RDB.SPopN(ctx, RedisTotalsDirty, 1000).Result()
	if err != nil && err != redis.Nil {
		connection.Log.Error("FlushDirtyTotals: SPopN failed", zap.Error(err))
		return
	}

	for _, m := range members {
		idcomp, idtrxStr, found := strings.Cut(m, "|")
		if !found {
			continue
		}
		idtrxkeluaran, err := strconv.Atoi(idtrxStr)
		if err != nil {
			continue
		}

		if err := repo.RefreshTotals(ctx, idcomp, idtrxkeluaran); err != nil {
			connection.Log.Error("FlushDirtyTotals: RefreshTotals failed",
				zap.Error(err), zap.String("idcomp", idcomp), zap.Int("idtrxkeluaran", idtrxkeluaran))
			// Transient DB error — re-mark dirty so this period isn't
			// silently dropped, instead of picked up again next tick.
			MarkTotalsDirty(idcomp, idtrxkeluaran)
			continue
		}

		// The totals just changed in Postgres — bust the agen dashboard's
		// cached copies now that there's something new to read. checkout.go
		// no longer busts these itself (it only marks the period dirty),
		// since doing so before this recompute ran would just repopulate the
		// cache with the same pre-checkout numbers.
		go connection.DeleteRedis(RedisTrxkeluaran + ":" + idcomp)
		go connection.DeleteRedis(RedisTrxkeluaranDetail + ":" + idcomp + ":" + strconv.Itoa(idtrxkeluaran))
	}
}
