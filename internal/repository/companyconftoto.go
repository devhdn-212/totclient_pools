package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/internal/config"
	"github.com/jackc/pgx/v5"
)

type companyconftotoRepository struct {
	db DBExecutor
}

func NewCompanyconftotoRepository(db DBExecutor) domain.CompanyconftotoRepository {
	return &companyconftotoRepository{
		db: db,
	}
}

func (c companyconftotoRepository) FindAll(ctx context.Context, idcompany string) ([]domain.Companyconftoto, error) {
	query := `SELECT * FROM ` + config.DB_tbl_companyconftoto + ` 
              WHERE idcompany = $1 
			  ORDER BY GREATEST(create_at, update_at) DESC`

	rows, err := c.db.Query(ctx, query, idcompany)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapping otomatis ke struct domain.Companyconftoto
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Companyconftoto])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c companyconftotoRepository) FindByID(ctx context.Context, idcompany string) (domain.Companyconftoto, error) {
	var compconftoto domain.Companyconftoto
	query := `SELECT idcompconftoto FROM ` + config.DB_tbl_companyconftoto + ` 
              WHERE idcompany = $1 LIMIT 1`

	err := c.db.QueryRow(ctx, query, idcompany).Scan(&compconftoto.IDcompconftoto)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Companyconftoto{}, nil
		}
		return compconftoto, err
	}
	return compconftoto, nil
}

func (c companyconftotoRepository) Save(ctx context.Context, compconftoto *domain.Companyconftoto) error {
	query := `INSERT INTO ` + config.DB_tbl_companyconftoto + ` 
                (
					idcompconftoto, idcompany, 
                	create_by, create_at
				) 
              	VALUES (
					$1, $2, $3, $4
				)`

	_, err := c.db.Exec(ctx, query,
		compconftoto.IDcompconftoto,
		compconftoto.IDcompany,
		compconftoto.CreateBy,
		compconftoto.CreateAt,
	)
	return err
}

func (c companyconftotoRepository) Update(ctx context.Context, compconftoto *domain.Companyconftoto) error {
	type col struct {
		name string
		val  any
	}

	cols := []col{
		{"angka_max_minbasket", compconftoto.AngkaMaxMinbasket},
		{"angka_max_minbet", compconftoto.AngkaMaxMinbet},
		{"angka_max_maxbet_4d", compconftoto.AngkaMaxMaxbet4d},
		{"angka_max_maxbet_3d", compconftoto.AngkaMaxMaxbet3d},
		{"angka_max_maxbet_3dd", compconftoto.AngkaMaxMaxbet3dd},
		{"angka_max_maxbet_2d", compconftoto.AngkaMaxMaxbet2d},
		{"angka_max_maxbet_2dd", compconftoto.AngkaMaxMaxbet2dd},
		{"angka_max_maxbet_2dt", compconftoto.AngkaMaxMaxbet2dt},
		{"angka_max_maxbet_4d_bbdisc", compconftoto.AngkaMaxMaxbet4dBbdisc},
		{"angka_max_maxbet_3d_bbdisc", compconftoto.AngkaMaxMaxbet3dBbdisc},
		{"angka_max_maxbet_3dd_bbdisc", compconftoto.AngkaMaxMaxbet3ddBbdisc},
		{"angka_max_maxbet_2d_bbdisc", compconftoto.AngkaMaxMaxbet2dBbdisc},
		{"angka_max_maxbet_2dd_bbdisc", compconftoto.AngkaMaxMaxbet2ddBbdisc},
		{"angka_max_maxbet_2dt_bbdisc", compconftoto.AngkaMaxMaxbet2dtBbdisc},
		{"angka_max_win4d_full", compconftoto.AngkaMaxWin4dFull},
		{"angka_max_win3d_full", compconftoto.AngkaMaxWin3dFull},
		{"angka_max_win3dd_full", compconftoto.AngkaMaxWin3ddFull},
		{"angka_max_win2d_full", compconftoto.AngkaMaxWin2dFull},
		{"angka_max_win2dd_full", compconftoto.AngkaMaxWin2ddFull},
		{"angka_max_win2dt_full", compconftoto.AngkaMaxWin2dtFull},
		{"angka_max_win4d_disc", compconftoto.AngkaMaxWin4dDisc},
		{"angka_max_win3d_disc", compconftoto.AngkaMaxWin3dDisc},
		{"angka_max_win3dd_disc", compconftoto.AngkaMaxWin3ddDisc},
		{"angka_max_win2d_disc", compconftoto.AngkaMaxWin2dDisc},
		{"angka_max_win2dd_disc", compconftoto.AngkaMaxWin2ddDisc},
		{"angka_max_win2dt_disc", compconftoto.AngkaMaxWin2dtDisc},
		{"angka_max_win4d_bb", compconftoto.AngkaMaxWin4dBb},
		{"angka_max_win3d_bb", compconftoto.AngkaMaxWin3dBb},
		{"angka_max_win3dd_bb", compconftoto.AngkaMaxWin3ddBb},
		{"angka_max_win2d_bb", compconftoto.AngkaMaxWin2dBb},
		{"angka_max_win2dd_bb", compconftoto.AngkaMaxWin2ddBb},
		{"angka_max_win2dt_bb", compconftoto.AngkaMaxWin2dtBb},
		{"angka_max_win4d_bb_kena", compconftoto.AngkaMaxWin4dBbKena},
		{"angka_max_win3d_bb_kena", compconftoto.AngkaMaxWin3dBbKena},
		{"angka_max_win3dd_bb_kena", compconftoto.AngkaMaxWin3ddBbKena},
		{"angka_max_win2d_bb_kena", compconftoto.AngkaMaxWin2dBbKena},
		{"angka_max_win2dd_bb_kena", compconftoto.AngkaMaxWin2ddBbKena},
		{"angka_max_win2dt_bb_kena", compconftoto.AngkaMaxWin2dtBbKena},
		{"cbebas_max_minbet", compconftoto.CbebasMaxMinbet},
		{"cbebas_max_maxbet", compconftoto.CbebasMaxMaxbet},
		{"cbebas_max_win", compconftoto.CbebasMaxWin},
		{"cmacau_max_minbet", compconftoto.CmacauMaxMinbet},
		{"cmacau_max_maxbet", compconftoto.CmacauMaxMaxbet},
		{"cmacau_max_win2", compconftoto.CmacauMaxWin2},
		{"cmacau_max_win3", compconftoto.CmacauMaxWin3},
		{"cmacau_max_win4", compconftoto.CmacauMaxWin4},
		{"cnaga_max_minbet", compconftoto.CnagaMaxMinbet},
		{"cnaga_max_maxbet", compconftoto.CnagaMaxMaxbet},
		{"cnaga_max_win3", compconftoto.CnagaMaxWin3},
		{"cnaga_max_win4", compconftoto.CnagaMaxWin4},
		{"cjitu_max_minbet", compconftoto.CjituMaxMinbet},
		{"cjitu_max_maxbet", compconftoto.CjituMaxMaxbet},
		{"cjitu_max_winas", compconftoto.CjituMaxWinas},
		{"cjitu_max_winkop", compconftoto.CjituMaxWinkop},
		{"cjitu_max_winkepala", compconftoto.CjituMaxWinkepala},
		{"cjitu_max_winekor", compconftoto.CjituMaxWinekor},
		{"umum50_max_minbet", compconftoto.Umum50MaxMinbet},
		{"umum50_max_maxbet", compconftoto.Umum50MaxMaxbet},
		{"special50_max_minbet", compconftoto.Special50MaxMinbet},
		{"special50_max_maxbet", compconftoto.Special50MaxMaxbet},
		{"kombinasi50_max_minbet", compconftoto.Kombinasi50MaxMinbet},
		{"kombinasi50_max_maxbet", compconftoto.Kombinasi50MaxMaxbet},
		{"macau_max_minbet", compconftoto.MacauMaxMinbet},
		{"macau_max_maxbet", compconftoto.MacauMaxMaxbet},
		{"macau_max_win", compconftoto.MacauMaxWin},
		{"dasar_max_minbet", compconftoto.DasarMaxMinbet},
		{"dasar_max_maxbet", compconftoto.DasarMaxMaxbet},
		{"shio_max_minbet", compconftoto.ShioMaxMinbet},
		{"shio_max_maxbet", compconftoto.ShioMaxMaxbet},
		{"shio_max_win", compconftoto.ShioMaxWin},
		{"shio_parent", compconftoto.ShioParent},
		{"update_by", compconftoto.UpdateBy},
		{"update_at", compconftoto.UpdateAt},
	}

	setClauses := make([]string, 0, len(cols))
	args := make([]any, 0, len(cols)+2)
	for i, cl := range cols {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", cl.name, i+1))
		args = append(args, cl.val)
	}

	n := len(cols)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE idcompconftoto = $%d AND idcompany = $%d`,
		config.DB_tbl_companyconftoto,
		strings.Join(setClauses, ", "),
		n+1, n+2,
	)
	args = append(args, compconftoto.IDcompconftoto, compconftoto.IDcompany)

	res, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
