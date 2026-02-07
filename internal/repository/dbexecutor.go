package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

/*
========================
LOW LEVEL SQL EXECUTOR
========================
*/

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SQLExecutor struct {
	DB *sql.DB
	Tx *sql.Tx
}

func NewSQLExecutor(db *sql.DB) DBExecutor {
	return &SQLExecutor{DB: db}
}

func NewSQLTxExecutor(tx *sql.Tx) DBExecutor {
	return &SQLExecutor{Tx: tx}
}

func (e *SQLExecutor) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	if e.Tx != nil {
		return e.Tx.ExecContext(ctx, q, args...)
	}
	return e.DB.ExecContext(ctx, q, args...)
}

func (e *SQLExecutor) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	if e.Tx != nil {
		return e.Tx.QueryContext(ctx, q, args...)
	}
	return e.DB.QueryContext(ctx, q, args...)
}

func (e *SQLExecutor) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	if e.Tx != nil {
		return e.Tx.QueryRowContext(ctx, q, args...)
	}
	return e.DB.QueryRowContext(ctx, q, args...)
}

/*
========================
GOQU DATASET PROVIDER
========================
*/

// ⬅️ INI YANG BENAR
type GoquDB interface {
	From(table ...interface{}) *goqu.SelectDataset
	Insert(table interface{}) *goqu.InsertDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
}

type GoquExecutor struct {
	DB   GoquDB
	Exec DBExecutor
}

func NewGoquExecutor(db *sql.DB) *GoquExecutor {
	return &GoquExecutor{
		DB:   goqu.New("default", db),
		Exec: NewSQLExecutor(db),
	}
}

func NewGoquTxExecutor(tx *sql.Tx) *GoquExecutor {
	return &GoquExecutor{
		DB:   goqu.NewTx("default", tx),
		Exec: NewSQLTxExecutor(tx),
	}
}
