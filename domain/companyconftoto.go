package domain

import (
	"context"
	"database/sql"

	"github.com/devhdn-212/totagen_api/dto"

	"github.com/shopspring/decimal"
)

type Companyconftoto struct {
	IDcompconftoto          string          `db:"idcompconftoto"`
	IDcompany               string          `db:"idcompany"`
	AngkaMaxMinbasket       decimal.Decimal `db:"angka_max_minbasket"`
	AngkaMaxMinbet          decimal.Decimal `db:"angka_max_minbet"`
	AngkaMaxMaxbet4d        decimal.Decimal `db:"angka_max_maxbet_4d"`
	AngkaMaxMaxbet3d        decimal.Decimal `db:"angka_max_maxbet_3d"`
	AngkaMaxMaxbet3dd       decimal.Decimal `db:"angka_max_maxbet_3dd"`
	AngkaMaxMaxbet2d        decimal.Decimal `db:"angka_max_maxbet_2d"`
	AngkaMaxMaxbet2dd       decimal.Decimal `db:"angka_max_maxbet_2dd"`
	AngkaMaxMaxbet2dt       decimal.Decimal `db:"angka_max_maxbet_2dt"`
	AngkaMaxMaxbet4dBbdisc  decimal.Decimal `db:"angka_max_maxbet_4d_bbdisc"`
	AngkaMaxMaxbet3dBbdisc  decimal.Decimal `db:"angka_max_maxbet_3d_bbdisc"`
	AngkaMaxMaxbet3ddBbdisc decimal.Decimal `db:"angka_max_maxbet_3dd_bbdisc"`
	AngkaMaxMaxbet2dBbdisc  decimal.Decimal `db:"angka_max_maxbet_2d_bbdisc"`
	AngkaMaxMaxbet2ddBbdisc decimal.Decimal `db:"angka_max_maxbet_2dd_bbdisc"`
	AngkaMaxMaxbet2dtBbdisc decimal.Decimal `db:"angka_max_maxbet_2dt_bbdisc"`
	AngkaMaxWin4dFull       decimal.Decimal `db:"angka_max_win4d_full"`
	AngkaMaxWin3dFull       decimal.Decimal `db:"angka_max_win3d_full"`
	AngkaMaxWin3ddFull      decimal.Decimal `db:"angka_max_win3dd_full"`
	AngkaMaxWin2dFull       decimal.Decimal `db:"angka_max_win2d_full"`
	AngkaMaxWin2ddFull      decimal.Decimal `db:"angka_max_win2dd_full"`
	AngkaMaxWin2dtFull      decimal.Decimal `db:"angka_max_win2dt_full"`
	AngkaMaxWin4dDisc       decimal.Decimal `db:"angka_max_win4d_disc"`
	AngkaMaxWin3dDisc       decimal.Decimal `db:"angka_max_win3d_disc"`
	AngkaMaxWin3ddDisc      decimal.Decimal `db:"angka_max_win3dd_disc"`
	AngkaMaxWin2dDisc       decimal.Decimal `db:"angka_max_win2d_disc"`
	AngkaMaxWin2ddDisc      decimal.Decimal `db:"angka_max_win2dd_disc"`
	AngkaMaxWin2dtDisc      decimal.Decimal `db:"angka_max_win2dt_disc"`
	AngkaMaxWin4dBb         decimal.Decimal `db:"angka_max_win4d_bb"`
	AngkaMaxWin3dBb         decimal.Decimal `db:"angka_max_win3d_bb"`
	AngkaMaxWin3ddBb        decimal.Decimal `db:"angka_max_win3dd_bb"`
	AngkaMaxWin2dBb         decimal.Decimal `db:"angka_max_win2d_bb"`
	AngkaMaxWin2ddBb        decimal.Decimal `db:"angka_max_win2dd_bb"`
	AngkaMaxWin2dtBb        decimal.Decimal `db:"angka_max_win2dt_bb"`
	AngkaMaxWin4dBbKena     decimal.Decimal `db:"angka_max_win4d_bb_kena"`
	AngkaMaxWin3dBbKena     decimal.Decimal `db:"angka_max_win3d_bb_kena"`
	AngkaMaxWin3ddBbKena    decimal.Decimal `db:"angka_max_win3dd_bb_kena"`
	AngkaMaxWin2dBbKena     decimal.Decimal `db:"angka_max_win2d_bb_kena"`
	AngkaMaxWin2ddBbKena    decimal.Decimal `db:"angka_max_win2dd_bb_kena"`
	AngkaMaxWin2dtBbKena    decimal.Decimal `db:"angka_max_win2dt_bb_kena"`
	CbebasMaxMinbet         decimal.Decimal `db:"cbebas_max_minbet"`
	CbebasMaxMaxbet         decimal.Decimal `db:"cbebas_max_maxbet"`
	CbebasMaxWin            decimal.Decimal `db:"cbebas_max_win"`
	CmacauMaxMinbet         decimal.Decimal `db:"cmacau_max_minbet"`
	CmacauMaxMaxbet         decimal.Decimal `db:"cmacau_max_maxbet"`
	CmacauMaxWin2           decimal.Decimal `db:"cmacau_max_win2"`
	CmacauMaxWin3           decimal.Decimal `db:"cmacau_max_win3"`
	CmacauMaxWin4           decimal.Decimal `db:"cmacau_max_win4"`
	CnagaMaxMinbet          decimal.Decimal `db:"cnaga_max_minbet"`
	CnagaMaxMaxbet          decimal.Decimal `db:"cnaga_max_maxbet"`
	CnagaMaxWin3            decimal.Decimal `db:"cnaga_max_win3"`
	CnagaMaxWin4            decimal.Decimal `db:"cnaga_max_win4"`
	CjituMaxMinbet          decimal.Decimal `db:"cjitu_max_minbet"`
	CjituMaxMaxbet          decimal.Decimal `db:"cjitu_max_maxbet"`
	CjituMaxWinas           decimal.Decimal `db:"cjitu_max_winas"`
	CjituMaxWinkop          decimal.Decimal `db:"cjitu_max_winkop"`
	CjituMaxWinkepala       decimal.Decimal `db:"cjitu_max_winkepala"`
	CjituMaxWinekor         decimal.Decimal `db:"cjitu_max_winekor"`
	Umum50MaxMinbet         decimal.Decimal `db:"umum50_max_minbet"`
	Umum50MaxMaxbet         decimal.Decimal `db:"umum50_max_maxbet"`
	Special50MaxMinbet      decimal.Decimal `db:"special50_max_minbet"`
	Special50MaxMaxbet      decimal.Decimal `db:"special50_max_maxbet"`
	Kombinasi50MaxMinbet    decimal.Decimal `db:"kombinasi50_max_minbet"`
	Kombinasi50MaxMaxbet    decimal.Decimal `db:"kombinasi50_max_maxbet"`
	MacauMaxMinbet          decimal.Decimal `db:"macau_max_minbet"`
	MacauMaxMaxbet          decimal.Decimal `db:"macau_max_maxbet"`
	MacauMaxWin             decimal.Decimal `db:"macau_max_win"`
	DasarMaxMinbet          decimal.Decimal `db:"dasar_max_minbet"`
	DasarMaxMaxbet          decimal.Decimal `db:"dasar_max_maxbet"`
	ShioMaxMinbet           decimal.Decimal `db:"shio_max_minbet"`
	ShioMaxMaxbet           decimal.Decimal `db:"shio_max_maxbet"`
	ShioMaxWin              decimal.Decimal `db:"shio_max_win"`
	ShioParent              int             `db:"shio_parent"`
	CreateBy                string          `db:"create_by"`
	CreateAt                sql.NullTime    `db:"create_at"`
	UpdateBy                string          `db:"update_by"`
	UpdateAt                sql.NullTime    `db:"update_at"`
}
type CompanyconftotoRepository interface {
	FindAll(ctx context.Context, idcomp string) ([]Companyconftoto, error)
	FindByID(ctx context.Context, idcomp string) (Companyconftoto, error)
	Save(ctx context.Context, compconftoto *Companyconftoto) error
	Update(ctx context.Context, compconftoto *Companyconftoto) error
}
type CompanyconftotoService interface {
	All(ctx context.Context, idcomp string) ([]dto.CompanyconftotoData, error)
	Save(ctx context.Context, req dto.CompanyconftotoSave, client string) error
}
