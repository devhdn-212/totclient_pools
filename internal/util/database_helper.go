package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetNextCounterManualTx(ctx context.Context, tx pgx.Tx, nmCounter string) (int64, error) {
	var currentCounter int64

	err := tx.QueryRow(ctx,
		`SELECT counter FROM tbl_counter WHERE nmcounter = $1 FOR UPDATE`,
		nmCounter,
	).Scan(&currentCounter)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("error select counter: %w", err)
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		finalCounter := currentCounter + 1
		_, err = tx.Exec(ctx,
			`UPDATE tbl_counter SET counter = $1 WHERE nmcounter = $2`,
			finalCounter, nmCounter,
		)
		if err != nil {
			return 0, fmt.Errorf("error update counter: %w", err)
		}
		return finalCounter, nil
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO tbl_counter (nmcounter, counter) VALUES ($1, $2)`,
		nmCounter, 1,
	)
	if err != nil {
		return 0, fmt.Errorf("error insert counter: %w", err)
	}
	return 1, nil
}
