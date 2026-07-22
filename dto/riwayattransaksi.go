package dto

type RiwayatTransaksiRequest struct {
	Token         string `json:"token" validate:"required"`
	Company       string `json:"company" validate:"required"`
	PasaranCode   string `json:"pasaran_code" validate:"required"`
	PasaranIdcomp string `json:"pasaran_idcomp" validate:"required"`
}

// RiwayatTransaksiResponse pairs each playerinvoice (one per checkout —
// the Transaksi menu's list view) with the individual bet rows that belong
// to it (the drill-down detail view, grouped client-side by playerinvoice).
type RiwayatTransaksiResponse struct {
	Members []TrxkeluaranmemberData `json:"members"`
	Details []TrxkeluarandetailData `json:"details"`
}
