package dto

type RiwayatTransaksiRequest struct {
	Token         string `json:"token" validate:"required"`
	Company       string `json:"company" validate:"required"`
	PasaranCode   string `json:"pasaran_code" validate:"required"`
	PasaranIdcomp string `json:"pasaran_idcomp" validate:"required"`
	// Idtrxkeluaran is optional. Omitted (0): only Periods is populated (the
	// period list). Set: Details is populated for that single period instead
	// — fetching every period's detail rows up front would be heavy, so the
	// client asks for one period's bets only when it actually opens it.
	Idtrxkeluaran int `json:"idtrxkeluaran,omitempty"`
}

// RiwayatTransaksiResponse is either the period list (Periods, one row per
// idtrxkeluaran) or one period's bet detail rows (Details) — never both;
// which one comes back depends on whether the request set Idtrxkeluaran.
type RiwayatTransaksiResponse struct {
	Periods []TrxkeluaranmemberPeriodData `json:"periods,omitempty"`
	Details []TrxkeluarandetailData       `json:"details,omitempty"`
}
