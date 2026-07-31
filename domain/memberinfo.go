package domain

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

// MothershipTransaction is one checkout chunk's ledger entry reported to the
// upstream member-wallet ("mothership") service — the authoritative debit of
// the player's balance. Exactly one of Credit/Debit must be greater than
// zero (mothership rejects otherwise); the consumer only ever sends Debit.
type MothershipTransaction struct {
	Idcompany     string
	Invoice       string
	Pasaran       string
	Playerinvoice string
	Username      string
	Credit        decimal.Decimal
	Debit         decimal.Decimal
}

type MothershipTransactionResult struct {
	Balance decimal.Decimal
	Status  string
}

// ErrInsufficientBalance is mothership's own live balance check on the
// debit call — more authoritative than whatever totclient_api validated
// against its (possibly stale-by-now) cached balance before publishing.
var ErrInsufficientBalance = errors.New("insufficient balance")

// MemberinfoService exposes the wallet debit plus a realtime balance
// re-check — unlike totclient_api, this worker never resolves a launch
// token into a session (that already happened before the checkout event
// was published), so it has no need for a full CheckToken method, just the
// raw balance lookup.
type MemberinfoService interface {
	// CheckBalance re-fetches the player's live wallet balance right before
	// debiting — balance is realtime and may have moved (other checkouts,
	// other games) since totclient_api's own balance check ran, which could
	// be anywhere from milliseconds to seconds earlier depending on Kafka
	// consumer lag.
	CheckBalance(ctx context.Context, idcompany, token string) (decimal.Decimal, error)
	SubmitTransaction(ctx context.Context, req MothershipTransaction) (*MothershipTransactionResult, error)
}
