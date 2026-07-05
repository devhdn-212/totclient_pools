package dto

import (
	"github.com/shopspring/decimal"
)

type CompanyconftotoData struct {
	IDcompconftoto          string          `json:"companyconftoto_idcompconftoto"`
	IDcompany               string          `json:"companyconftoto_idcompany"`
	AngkaMaxMinbasket       decimal.Decimal `json:"companyconftoto_angka_max_minbasket"`
	AngkaMaxMinbet          decimal.Decimal `json:"companyconftoto_angka_max_minbet"`
	AngkaMaxMaxbet4d        decimal.Decimal `json:"companyconftoto_angka_max_maxbet_4d"`
	AngkaMaxMaxbet3d        decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3d"`
	AngkaMaxMaxbet3dd       decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3dd"`
	AngkaMaxMaxbet2d        decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2d"`
	AngkaMaxMaxbet2dd       decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dd"`
	AngkaMaxMaxbet2dt       decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dt"`
	AngkaMaxMaxbet4dBbdisc  decimal.Decimal `json:"companyconftoto_angka_max_maxbet_4d_bbdisc"`
	AngkaMaxMaxbet3dBbdisc  decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3d_bbdisc"`
	AngkaMaxMaxbet3ddBbdisc decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3dd_bbdisc"`
	AngkaMaxMaxbet2dBbdisc  decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2d_bbdisc"`
	AngkaMaxMaxbet2ddBbdisc decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dd_bbdisc"`
	AngkaMaxMaxbet2dtBbdisc decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dt_bbdisc"`
	AngkaMaxWin4dFull       decimal.Decimal `json:"companyconftoto_angka_max_win4d_full"`
	AngkaMaxWin3dFull       decimal.Decimal `json:"companyconftoto_angka_max_win3d_full"`
	AngkaMaxWin3ddFull      decimal.Decimal `json:"companyconftoto_angka_max_win3dd_full"`
	AngkaMaxWin2dFull       decimal.Decimal `json:"companyconftoto_angka_max_win2d_full"`
	AngkaMaxWin2ddFull      decimal.Decimal `json:"companyconftoto_angka_max_win2dd_full"`
	AngkaMaxWin2dtFull      decimal.Decimal `json:"companyconftoto_angka_max_win2dt_full"`
	AngkaMaxWin4dDisc       decimal.Decimal `json:"companyconftoto_angka_max_win4d_disc"`
	AngkaMaxWin3dDisc       decimal.Decimal `json:"companyconftoto_angka_max_win3d_disc"`
	AngkaMaxWin3ddDisc      decimal.Decimal `json:"companyconftoto_angka_max_win3dd_disc"`
	AngkaMaxWin2dDisc       decimal.Decimal `json:"companyconftoto_angka_max_win2d_disc"`
	AngkaMaxWin2ddDisc      decimal.Decimal `json:"companyconftoto_angka_max_win2dd_disc"`
	AngkaMaxWin2dtDisc      decimal.Decimal `json:"companyconftoto_angka_max_win2dt_disc"`
	AngkaMaxWin4dBb         decimal.Decimal `json:"companyconftoto_angka_max_win4d_bb"`
	AngkaMaxWin3dBb         decimal.Decimal `json:"companyconftoto_angka_max_win3d_bb"`
	AngkaMaxWin3ddBb        decimal.Decimal `json:"companyconftoto_angka_max_win3dd_bb"`
	AngkaMaxWin2dBb         decimal.Decimal `json:"companyconftoto_angka_max_win2d_bb"`
	AngkaMaxWin2ddBb        decimal.Decimal `json:"companyconftoto_angka_max_win2dd_bb"`
	AngkaMaxWin2dtBb        decimal.Decimal `json:"companyconftoto_angka_max_win2dt_bb"`
	AngkaMaxWin4dBbKena     decimal.Decimal `json:"companyconftoto_angka_max_win4d_bb_kena"`
	AngkaMaxWin3dBbKena     decimal.Decimal `json:"companyconftoto_angka_max_win3d_bb_kena"`
	AngkaMaxWin3ddBbKena    decimal.Decimal `json:"companyconftoto_angka_max_win3dd_bb_kena"`
	AngkaMaxWin2dBbKena     decimal.Decimal `json:"companyconftoto_angka_max_win2d_bb_kena"`
	AngkaMaxWin2ddBbKena    decimal.Decimal `json:"companyconftoto_angka_max_win2dd_bb_kena"`
	AngkaMaxWin2dtBbKena    decimal.Decimal `json:"companyconftoto_angka_max_win2dt_bb_kena"`
	CbebasMaxMinbet         decimal.Decimal `json:"companyconftoto_cbebas_max_minbet"`
	CbebasMaxMaxbet         decimal.Decimal `json:"companyconftoto_cbebas_max_maxbet"`
	CbebasMaxWin            decimal.Decimal `json:"companyconftoto_cbebas_max_win"`
	CmacauMaxMinbet         decimal.Decimal `json:"companyconftoto_cmacau_max_minbet"`
	CmacauMaxMaxbet         decimal.Decimal `json:"companyconftoto_cmacau_max_maxbet"`
	CmacauMaxWin2           decimal.Decimal `json:"companyconftoto_cmacau_max_win2"`
	CmacauMaxWin3           decimal.Decimal `json:"companyconftoto_cmacau_max_win3"`
	CmacauMaxWin4           decimal.Decimal `json:"companyconftoto_cmacau_max_win4"`
	CnagaMaxMinbet          decimal.Decimal `json:"companyconftoto_cnaga_max_minbet"`
	CnagaMaxMaxbet          decimal.Decimal `json:"companyconftoto_cnaga_max_maxbet"`
	CnagaMaxWin3            decimal.Decimal `json:"companyconftoto_cnaga_max_win3"`
	CnagaMaxWin4            decimal.Decimal `json:"companyconftoto_cnaga_max_win4"`
	CjituMaxMinbet          decimal.Decimal `json:"companyconftoto_cjitu_max_minbet"`
	CjituMaxMaxbet          decimal.Decimal `json:"companyconftoto_cjitu_max_maxbet"`
	CjituMaxWinas           decimal.Decimal `json:"companyconftoto_cjitu_max_winas"`
	CjituMaxWinkop          decimal.Decimal `json:"companyconftoto_cjitu_max_winkop"`
	CjituMaxWinkepala       decimal.Decimal `json:"companyconftoto_cjitu_max_winkepala"`
	CjituMaxWinekor         decimal.Decimal `json:"companyconftoto_cjitu_max_winekor"`
	Umum50MaxMinbet         decimal.Decimal `json:"companyconftoto_umum50_max_minbet"`
	Umum50MaxMaxbet         decimal.Decimal `json:"companyconftoto_umum50_max_maxbet"`
	Special50MaxMinbet      decimal.Decimal `json:"companyconftoto_special50_max_minbet"`
	Special50MaxMaxbet      decimal.Decimal `json:"companyconftoto_special50_max_maxbet"`
	Kombinasi50MaxMinbet    decimal.Decimal `json:"companyconftoto_kombinasi50_max_minbet"`
	Kombinasi50MaxMaxbet    decimal.Decimal `json:"companyconftoto_kombinasi50_max_maxbet"`
	MacauMaxMinbet          decimal.Decimal `json:"companyconftoto_macau_max_minbet"`
	MacauMaxMaxbet          decimal.Decimal `json:"companyconftoto_macau_max_maxbet"`
	MacauMaxWin             decimal.Decimal `json:"companyconftoto_macau_max_win"`
	DasarMaxMinbet          decimal.Decimal `json:"companyconftoto_dasar_max_minbet"`
	DasarMaxMaxbet          decimal.Decimal `json:"companyconftoto_dasar_max_maxbet"`
	ShioMaxMinbet           decimal.Decimal `json:"companyconftoto_shio_max_minbet"`
	ShioMaxMaxbet           decimal.Decimal `json:"companyconftoto_shio_max_maxbet"`
	ShioMaxWin              decimal.Decimal `json:"companyconftoto_shio_max_win"`
	ShioParent              int             `json:"companyconftoto_shio_parent"`
	CreateBy                string          `json:"companyconftoto_create_by"`
}
type CompanyconftotoAll struct {
	IDcompany string `json:"compconfoto_idcompany" validate:"required"`
}
type CompanyconftotoSave struct {
	Type                    string          `json:"type" validate:"required"`
	IDcompconftoto          string          `json:"companyconftoto_idcompconftoto" validate:"required"`
	IDcompany               string          `json:"companyconftoto_idcompany"  validate:"required"`
	AngkaMaxMinbasket       decimal.Decimal `json:"companyconftoto_angka_max_minbasket"`
	AngkaMaxMinbet          decimal.Decimal `json:"companyconftoto_angka_max_minbet"`
	AngkaMaxMaxbet4d        decimal.Decimal `json:"companyconftoto_angka_max_maxbet_4d"`
	AngkaMaxMaxbet3d        decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3d"`
	AngkaMaxMaxbet3dd       decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3dd"`
	AngkaMaxMaxbet2d        decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2d"`
	AngkaMaxMaxbet2dd       decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dd"`
	AngkaMaxMaxbet2dt       decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dt"`
	AngkaMaxMaxbet4dBbdisc  decimal.Decimal `json:"companyconftoto_angka_max_maxbet_4d_bbdisc"`
	AngkaMaxMaxbet3dBbdisc  decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3d_bbdisc"`
	AngkaMaxMaxbet3ddBbdisc decimal.Decimal `json:"companyconftoto_angka_max_maxbet_3dd_bbdisc"`
	AngkaMaxMaxbet2dBbdisc  decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2d_bbdisc"`
	AngkaMaxMaxbet2ddBbdisc decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dd_bbdisc"`
	AngkaMaxMaxbet2dtBbdisc decimal.Decimal `json:"companyconftoto_angka_max_maxbet_2dt_bbdisc"`
	AngkaMaxWin4dFull       decimal.Decimal `json:"companyconftoto_angka_max_win4d_full"`
	AngkaMaxWin3dFull       decimal.Decimal `json:"companyconftoto_angka_max_win3d_full"`
	AngkaMaxWin3ddFull      decimal.Decimal `json:"companyconftoto_angka_max_win3dd_full"`
	AngkaMaxWin2dFull       decimal.Decimal `json:"companyconftoto_angka_max_win2d_full"`
	AngkaMaxWin2ddFull      decimal.Decimal `json:"companyconftoto_angka_max_win2dd_full"`
	AngkaMaxWin2dtFull      decimal.Decimal `json:"companyconftoto_angka_max_win2dt_full"`
	AngkaMaxWin4dDisc       decimal.Decimal `json:"companyconftoto_angka_max_win4d_disc"`
	AngkaMaxWin3dDisc       decimal.Decimal `json:"companyconftoto_angka_max_win3d_disc"`
	AngkaMaxWin3ddDisc      decimal.Decimal `json:"companyconftoto_angka_max_win3dd_disc"`
	AngkaMaxWin2dDisc       decimal.Decimal `json:"companyconftoto_angka_max_win2d_disc"`
	AngkaMaxWin2ddDisc      decimal.Decimal `json:"companyconftoto_angka_max_win2dd_disc"`
	AngkaMaxWin2dtDisc      decimal.Decimal `json:"companyconftoto_angka_max_win2dt_disc"`
	AngkaMaxWin4dBb         decimal.Decimal `json:"companyconftoto_angka_max_win4d_bb"`
	AngkaMaxWin3dBb         decimal.Decimal `json:"companyconftoto_angka_max_win3d_bb"`
	AngkaMaxWin3ddBb        decimal.Decimal `json:"companyconftoto_angka_max_win3dd_bb"`
	AngkaMaxWin2dBb         decimal.Decimal `json:"companyconftoto_angka_max_win2d_bb"`
	AngkaMaxWin2ddBb        decimal.Decimal `json:"companyconftoto_angka_max_win2dd_bb"`
	AngkaMaxWin2dtBb        decimal.Decimal `json:"companyconftoto_angka_max_win2dt_bb"`
	AngkaMaxWin4dBbKena     decimal.Decimal `json:"companyconftoto_angka_max_win4d_bb_kena"`
	AngkaMaxWin3dBbKena     decimal.Decimal `json:"companyconftoto_angka_max_win3d_bb_kena"`
	AngkaMaxWin3ddBbKena    decimal.Decimal `json:"companyconftoto_angka_max_win3dd_bb_kena"`
	AngkaMaxWin2dBbKena     decimal.Decimal `json:"companyconftoto_angka_max_win2d_bb_kena"`
	AngkaMaxWin2ddBbKena    decimal.Decimal `json:"companyconftoto_angka_max_win2dd_bb_kena"`
	AngkaMaxWin2dtBbKena    decimal.Decimal `json:"companyconftoto_angka_max_win2dt_bb_kena"`
	CbebasMaxMinbet         decimal.Decimal `json:"companyconftoto_cbebas_max_minbet"`
	CbebasMaxMaxbet         decimal.Decimal `json:"companyconftoto_cbebas_max_maxbet"`
	CbebasMaxWin            decimal.Decimal `json:"companyconftoto_cbebas_max_win"`
	CmacauMaxMinbet         decimal.Decimal `json:"companyconftoto_cmacau_max_minbet"`
	CmacauMaxMaxbet         decimal.Decimal `json:"companyconftoto_cmacau_max_maxbet"`
	CmacauMaxWin2           decimal.Decimal `json:"companyconftoto_cmacau_max_win2"`
	CmacauMaxWin3           decimal.Decimal `json:"companyconftoto_cmacau_max_win3"`
	CmacauMaxWin4           decimal.Decimal `json:"companyconftoto_cmacau_max_win4"`
	CnagaMaxMinbet          decimal.Decimal `json:"companyconftoto_cnaga_max_minbet"`
	CnagaMaxMaxbet          decimal.Decimal `json:"companyconftoto_cnaga_max_maxbet"`
	CnagaMaxWin3            decimal.Decimal `json:"companyconftoto_cnaga_max_win3"`
	CnagaMaxWin4            decimal.Decimal `json:"companyconftoto_cnaga_max_win4"`
	CjituMaxMinbet          decimal.Decimal `json:"companyconftoto_cjitu_max_minbet"`
	CjituMaxMaxbet          decimal.Decimal `json:"companyconftoto_cjitu_max_maxbet"`
	CjituMaxWinas           decimal.Decimal `json:"companyconftoto_cjitu_max_winas"`
	CjituMaxWinkop          decimal.Decimal `json:"companyconftoto_cjitu_max_winkop"`
	CjituMaxWinkepala       decimal.Decimal `json:"companyconftoto_cjitu_max_winkepala"`
	CjituMaxWinekor         decimal.Decimal `json:"companyconftoto_cjitu_max_winekor"`
	Umum50MaxMinbet         decimal.Decimal `json:"companyconftoto_umum50_max_minbet"`
	Umum50MaxMaxbet         decimal.Decimal `json:"companyconftoto_umum50_max_maxbet"`
	Special50MaxMinbet      decimal.Decimal `json:"companyconftoto_special50_max_minbet"`
	Special50MaxMaxbet      decimal.Decimal `json:"companyconftoto_special50_max_maxbet"`
	Kombinasi50MaxMinbet    decimal.Decimal `json:"companyconftoto_kombinasi50_max_minbet"`
	Kombinasi50MaxMaxbet    decimal.Decimal `json:"companyconftoto_kombinasi50_max_maxbet"`
	MacauMaxMinbet          decimal.Decimal `json:"companyconftoto_macau_max_minbet"`
	MacauMaxMaxbet          decimal.Decimal `json:"companyconftoto_macau_max_maxbet"`
	MacauMaxWin             decimal.Decimal `json:"companyconftoto_macau_max_win"`
	DasarMaxMinbet          decimal.Decimal `json:"companyconftoto_dasar_max_minbet"`
	DasarMaxMaxbet          decimal.Decimal `json:"companyconftoto_dasar_max_maxbet"`
	ShioMaxMinbet           decimal.Decimal `json:"companyconftoto_shio_max_minbet"`
	ShioMaxMaxbet           decimal.Decimal `json:"companyconftoto_shio_max_maxbet"`
	ShioMaxWin              decimal.Decimal `json:"companyconftoto_shio_max_win"`
	ShioParent              int             `json:"companyconftoto_shio_parent"`
}
