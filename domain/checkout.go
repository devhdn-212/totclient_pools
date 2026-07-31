package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// CheckoutKafkaItem is one accepted bet line within a CheckoutKafkaMessage —
// mirrors totclient_api's domain.CheckoutKafkaItem field-for-field (two
// separate Go modules, so the shape is defined again here rather than
// shared — same pattern already used between totclient_client and
// totclient_api, see DOKUMENTASI.md).
type CheckoutKafkaItem struct {
	DetailID    string          `json:"detail_id"`
	Nomor       string          `json:"nomor"`
	Typegame    string          `json:"typegame"`
	Posisitogel string          `json:"posisitogel"`
	Bet         int             `json:"bet"`
	DiscRate    decimal.Decimal `json:"disc_rate"`
	KeiRate     decimal.Decimal `json:"kei_rate"`
	Win         decimal.Decimal `json:"win"`
	Payout      decimal.Decimal `json:"payout"`
}

// CheckoutKafkaMessage is the checkout event consumed off Kafka — see
// totclient_api's domain/checkout.go for the producer side and the full
// field-by-field rationale. totclient_api has already done every validation
// (balance, limittotal/limitglobal, pasaran open/closed, item format/range)
// before publishing this, so the consumer trusts it completely: debit the
// wallet, then persist exactly what's here.
type CheckoutKafkaMessage struct {
	Idcompany        string `json:"idcompany"`
	IDtrxkeluaran    int    `json:"idtrxkeluaran"`
	Invoice          string `json:"invoice"`
	Pasaran          string `json:"pasaran"`
	PasaranIdcomp    string `json:"pasaran_idcomp"`
	Playerinvoice    int    `json:"playerinvoice"`
	PlayerinvoiceStr string `json:"playerinvoice_str"`
	Betround         int    `json:"betround"`
	Username         string `json:"username"`
	// Token is the player's launch token — used to re-check the live wallet
	// balance right before debiting, since balance is realtime and may have
	// moved between totclient_api's own balance check and whenever this
	// event actually gets consumed here.
	Token              string              `json:"token"`
	Ipaddress          string              `json:"ipaddress"`
	Devicemember       string              `json:"devicemember"`
	Datekeluarandetail time.Time           `json:"datekeluarandetail"`
	MemberID           string              `json:"member_id"`
	Totalbet           decimal.Decimal     `json:"totalbet"`
	Totalbayar         decimal.Decimal     `json:"totalbayar"`
	Totaldiscount      decimal.Decimal     `json:"totaldiscount"`
	Totalkei           decimal.Decimal     `json:"totalkei"`
	Totalpair          int                 `json:"totalpair"`
	Items              []CheckoutKafkaItem `json:"items"`
}
