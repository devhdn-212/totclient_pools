package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

// Status values for MemberInvoice.Status — must match tbl_trx_member_invoice's
// chk_invoice_status CHECK constraint exactly (title case, not upper/lower).
const (
	MemberInvoiceRequested = "Requested"
	MemberInvoiceCompleted = "Completed"
	MemberInvoiceVoid      = "Void"
)

// MemberInvoice mirrors public.tbl_trx_member_invoice — a cross-company
// reconciliation table (schema public, not per-company db_tot_<company>)
// tracking every checkout's wallet-debit attempt independently of whether
// trxkeluarandetail/trxkeluaranmember ever actually got persisted. See
// DOKUMENTASI.md § 9.8 for the full lifecycle/rationale.
//
// PK is (agent_code, invoice_id, date_transaction) — invoice_id
// (playerinvoice) is only unique WITHIN one company's own idtrxkeluaran
// sequence, so agent_code has to be part of the key or two different
// companies could collide on the same invoice_id+date.
type MemberInvoice struct {
	ID              string          `db:"id"`
	AgentCode       string          `db:"agent_code"`
	InvoiceID       int             `db:"invoice_id"`
	Username        string          `db:"username"`
	PlayerToken     string          `db:"player_token"`
	DebitAmount     decimal.Decimal `db:"debit_amount"`
	DateTransaction time.Time       `db:"date_transaction"`
	Status          string          `db:"status"`
	CreatedAt       sql.NullTime    `db:"create_at"`
	UpdatedAt       sql.NullTime    `db:"update_at"`
}

// MemberInvoiceRepository is the totclient_pools-side write path: insert the
// Requested row right before the wallet debit call, then transition it to
// Completed once persist() succeeds in the same tx. Rows left at Requested -
// whether the wallet debit itself failed, OR it succeeded but persist()
// didn't - are BOTH deliberately left alone here (no Void transition on this
// side at all, dikonfirmasi user): totagen_api's "Pending Bet Request" panel
// is what reads them and is the only place that ever moves Requested to
// Void (manual admin refund, reason typically "Failed to deduct" — see that
// repo's domain/memberinvoice.go and documentasi.md bagian 64).
type MemberInvoiceRepository interface {
	InsertRequested(ctx context.Context, inv *MemberInvoice) error
	MarkCompleted(ctx context.Context, agentCode string, invoiceID int, dateTransaction time.Time) error
	EnsureMonthPartition(ctx context.Context, monthStart time.Time) error
}
