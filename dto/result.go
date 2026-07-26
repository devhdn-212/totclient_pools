package dto

type ResultRequest struct {
	Token         string `json:"token" validate:"required"`
	Company       string `json:"company" validate:"required"`
	PasaranCode   string `json:"pasaran_code" validate:"required"`
	PasaranIdcomp string `json:"pasaran_idcomp" validate:"required"`
	// Bulan is the month to show results for (1-12) — the year is always
	// the current year, never client-controlled. Omitted or out of range
	// defaults to the current month.
	Bulan int `json:"bulan,omitempty"`
}

// ResultItemData is one decided period ("Result" menu) — public draw-result
// information for the pasaran, not scoped to any one player.
type ResultItemData struct {
	Periode          string `json:"periode"`
	Datekeluaran     string `json:"datekeluaran"`
	Keluarantogel    string `json:"keluarantogel"`
	Aliascomppasaran string `json:"aliascomppasaran"`
}

type ResultResponse struct {
	Results []ResultItemData `json:"results"`
}
