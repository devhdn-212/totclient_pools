package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
========================
LOW LEVEL PGX EXECUTOR
========================
*/

// DBExecutor disesuaikan dengan signature metode milik pgx
type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type PGXExecutor struct {
	// Pool digunakan untuk koneksi standar, Tx untuk transaksi
	Pool *pgxpool.Pool
	Tx   pgx.Tx
}

// NewPGXExecutor inisialisasi dengan pgxpool
func NewPGXExecutor(pool *pgxpool.Pool) DBExecutor {
	return &PGXExecutor{Pool: pool}
}

// NewPGXTxExecutor inisialisasi dengan transaksi pgx
func NewPGXTxExecutor(tx pgx.Tx) DBExecutor {
	return &PGXExecutor{Tx: tx}
}

func (e *PGXExecutor) Exec(ctx context.Context, q string, args ...any) (pgconn.CommandTag, error) {
	if e.Tx != nil {
		return e.Tx.Exec(ctx, q, args...)
	}
	return e.Pool.Exec(ctx, q, args...)
}

func (e *PGXExecutor) Query(ctx context.Context, q string, args ...any) (pgx.Rows, error) {
	if e.Tx != nil {
		return e.Tx.Query(ctx, q, args...)
	}
	return e.Pool.Query(ctx, q, args...)
}

func (e *PGXExecutor) QueryRow(ctx context.Context, q string, args ...any) pgx.Row {
	if e.Tx != nil {
		return e.Tx.QueryRow(ctx, q, args...)
	}
	return e.Pool.QueryRow(ctx, q, args...)
}

/*
========================
REPOSITORI WRAPPER
========================
*/

// Karena goqu dihapus, kita tidak perlu GoquDB interface lagi.
// Kita buat struct sederhana untuk membungkus executor.
type BaseRepository struct {
	DB DBExecutor
}

func NewRepository(db *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{
		DB: NewPGXExecutor(db),
	}
}

func NewTxRepository(tx pgx.Tx) *BaseRepository {
	return &BaseRepository{
		DB: NewPGXTxExecutor(tx),
	}
}
