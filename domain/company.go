package domain

import "context"

type Company struct {
	IDcompany   string `db:"idcompany"`
	Compname    string `db:"compname"`
	Compstatus  string `db:"compstatus"`
}

type CompanyRepository interface {
	// FindByID looks up an agent by its idcompany code. Returns a zero-value
	// Company (IDcompany == "") when no active row matches, same not-found
	// convention as PasaranRepository.FindByID.
	FindByID(ctx context.Context, idcompany string) (Company, error)
}
