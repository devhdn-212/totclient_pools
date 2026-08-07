package repository

import (
	"context"
	"time"

	"github.com/devhdn-212/totclient_pools/domain"
)

// memberInvoiceTable is schema-less on purpose — public.tbl_trx_member_invoice
// is shared across every company (unlike the per-company db_tot_<company>
// tables elsewhere in this repo, see util.Get_mapping_totodb), so there's no
// idcomp to resolve a schema from here.
const memberInvoiceTable = "public.tbl_trx_member_invoice"

type memberInvoiceRepository struct {
	db DBExecutor
}

func NewMemberInvoiceRepository(db DBExecutor) domain.MemberInvoiceRepository {
	return &memberInvoiceRepository{db: db}
}

func (a memberInvoiceRepository) InsertRequested(ctx context.Context, inv *domain.MemberInvoice) error {
	query := `INSERT INTO ` + memberInvoiceTable + `
		(agent_code, invoice_id, username, player_token, debit_amount, date_transaction, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := a.db.Exec(ctx, query,
		inv.AgentCode, inv.InvoiceID, inv.Username, inv.PlayerToken, inv.DebitAmount, inv.DateTransaction, domain.MemberInvoiceRequested,
	)
	return err
}

// MarkCompleted only moves a row OUT of Requested — WHERE status = 'Requested'
// guards against updating a row some other caller already resolved (belt and
// suspenders; in practice each invoice's status is only ever touched once by
// the checkout flow that created it).
func (a memberInvoiceRepository) MarkCompleted(ctx context.Context, agentCode string, invoiceID int, dateTransaction time.Time) error {
	query := `UPDATE ` + memberInvoiceTable + `
		SET status = $1, update_at = now()
		WHERE agent_code = $2 AND invoice_id = $3 AND date_transaction = $4 AND status = $5`
	_, err := a.db.Exec(ctx, query, domain.MemberInvoiceCompleted, agentCode, invoiceID, dateTransaction, domain.MemberInvoiceRequested)
	return err
}

// EnsureMonthPartition mirrors totagen_api/totagen_pools's ensureMonthPartitions
// (see that repo's documentasi.md bagian 63) but for this one public-schema
// table instead of a per-company one - idempotent via "CREATE TABLE IF NOT
// EXISTS ... PARTITION OF ...", safe to call repeatedly/concurrently.
func (a memberInvoiceRepository) EnsureMonthPartition(ctx context.Context, monthStart time.Time) error {
	monthEnd := monthStart.AddDate(0, 1, 0)
	suffix := monthStart.Format("200601")
	partitionName := "public.tbl_trx_member_invoice_p" + suffix

	query := `CREATE TABLE IF NOT EXISTS ` + partitionName + `
		PARTITION OF ` + memberInvoiceTable + `
		FOR VALUES FROM ('` + monthStart.Format("2006-01-02") + `') TO ('` + monthEnd.Format("2006-01-02") + `')`
	_, err := a.db.Exec(ctx, query)
	return err
}
