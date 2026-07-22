package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// CheckAndIncrementLimitTx atomically adds amount to a running total kept in
// tbl_counter (same table GetNextCounterManualTx uses, different key
// namespace) and only commits the increment if the new total stays within
// limit. limit <= 0 means uncapped — always succeeds. Locks only the one
// counter row (SELECT ... FOR UPDATE), so unrelated keys never block each
// other; concurrent bets on the SAME key are serialized, which is exactly
// what's needed to keep limittotal/limitglobal correct under concurrency.
func CheckAndIncrementLimitTx(ctx context.Context, tx pgx.Tx, nmCounter string, amount, limit int64) (newTotal int64, ok bool, err error) {
	if limit <= 0 {
		return 0, true, nil
	}

	var current int64
	err = tx.QueryRow(ctx,
		`SELECT counter FROM tbl_counter WHERE nmcounter = $1 FOR UPDATE`,
		nmCounter,
	).Scan(&current)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, fmt.Errorf("error select limit counter: %w", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		current = 0
		if _, err = tx.Exec(ctx,
			`INSERT INTO tbl_counter (nmcounter, counter) VALUES ($1, 0)`,
			nmCounter,
		); err != nil {
			return 0, false, fmt.Errorf("error insert limit counter: %w", err)
		}
	}

	if current+amount > limit {
		return current, false, nil
	}

	newTotal = current + amount
	if _, err = tx.Exec(ctx,
		`UPDATE tbl_counter SET counter = $1 WHERE nmcounter = $2`,
		newTotal, nmCounter,
	); err != nil {
		return 0, false, fmt.Errorf("error update limit counter: %w", err)
	}
	return newTotal, true, nil
}

// DecrementCounterTx undoes a prior CheckAndIncrementLimitTx increment —
// used when a bet passes the limittotal check but then fails limitglobal,
// so the limittotal counter it already bumped needs to be rolled back
// within the same transaction (the bet as a whole is rejected).
func DecrementCounterTx(ctx context.Context, tx pgx.Tx, nmCounter string, amount int64) error {
	_, err := tx.Exec(ctx,
		`UPDATE tbl_counter SET counter = counter - $1 WHERE nmcounter = $2`,
		amount, nmCounter,
	)
	if err != nil {
		return fmt.Errorf("error decrement limit counter: %w", err)
	}
	return nil
}
