package dto

import "github.com/shopspring/decimal"

// CheckoutItem mirrors one entry of the client's `data: chunk` payload — the
// client-computed diskon/win/kei/bayar fields are only trusted for display
// on the client; the server recomputes/validates them independently.
type CheckoutItem struct {
	ID           string          `json:"id" validate:"required"`
	Nomor        string          `json:"nomor" validate:"required"`
	Permainan    string          `json:"permainan" validate:"required"`
	Bet          int             `json:"bet" validate:"required,gt=0"`
	Diskon       decimal.Decimal `json:"diskon"`
	Diskonpercen decimal.Decimal `json:"diskonpercen"`
	Bayar        decimal.Decimal `json:"bayar"`
	Win          decimal.Decimal `json:"win"`
	Kei          decimal.Decimal `json:"kei"`
	KeiPercen    decimal.Decimal `json:"kei_percen"`
	Tipetoto     string          `json:"tipetoto" validate:"required"`
}

type CheckoutRequest struct {
	PasaranIdtransaction int             `json:"pasaran_idtransaction" validate:"required"`
	PasaranIdcomp        string          `json:"pasaran_idcomp" validate:"required"`
	PasaranCode          string          `json:"pasaran_code" validate:"required"`
	PasaranPeriode       int             `json:"pasaran_periode"`
	Token                string          `json:"token" validate:"required"`
	Company              string          `json:"company" validate:"required"`
	Username             string          `json:"username"`
	Ipaddress            string          `json:"ipaddress"`
	Devicemember         string          `json:"devicemember"`
	Timezone             string          `json:"timezone"`
	Data                 []CheckoutItem  `json:"data" validate:"required,dive"`
	Total                decimal.Decimal `json:"total"`
}

// CheckoutItemResult reports what happened to one submitted item —
// "accepted" or "rejected" (with why, and how much headroom is left on
// whichever limit it tripped, so the player can decide to resubmit lower).
type CheckoutItemResult struct {
	ID        string          `json:"id"`
	Status    string          `json:"status"`
	Reason    string          `json:"reason,omitempty"`
	Sisalimit decimal.Decimal `json:"sisalimit,omitempty"`
}

type CheckoutResponse struct {
	Playerinvoice int                  `json:"playerinvoice"`
	Totalbayar    decimal.Decimal      `json:"totalbayar"`
	Items         []CheckoutItemResult `json:"items"`
}
