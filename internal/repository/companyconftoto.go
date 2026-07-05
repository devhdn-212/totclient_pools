package repository

import (
	"context"
	"errors"

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
              ORDER BY create_at, update_at DESC`

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

func (c companyconftotoRepository) FindByID(ctx context.Context, idcompany, id string) (domain.Companyconftoto, error) {
	var compconftoto domain.Companyconftoto
	query := `SELECT idcompconftoto FROM ` + config.DB_tbl_companyconftoto + ` 
              WHERE idcompconftoto = $1 AND idcompany = $2 LIMIT 1`

	err := c.db.QueryRow(ctx, query, id, idcompany).Scan(&compconftoto.IDcompconftoto, &compconftoto.IDcompany)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Companyconftoto{}, nil
		}
		return compconftoto, err
	}
	return compconftoto, nil
}

func (c companyconftotoRepository) SaveAngka(ctx context.Context, compconftoto *domain.Companyconftoto) error {
	query := `INSERT INTO ` + config.DB_tbl_companyconftoto + ` 
                (
					idcompconftoto, idcompany, 
					angka_max_minbasket, angka_max_minbet, 
					angka_max_maxbet_4d, 
					angka_max_maxbet_3d, angka_max_maxbet_3dd, 
					angka_max_maxbet_2d, angka_max_maxbet_2dd, angka_max_maxbet_2dt, 
					angka_max_maxbet_4d_bbdisc, 
					angka_max_maxbet_3d_bbdisc, angka_max_maxbet_3dd_bbdisc,  
					angka_max_maxbet_2d_bbdisc, angka_max_maxbet_2dd_bbdisc, angka_max_maxbet_2dt_bbdisc, 
					angka_max_win4d_full, 
					angka_max_win3d_full, angka_max_win3dd_full, 
					angka_max_win2d_full,angka_max_win2dd_full,angka_max_win2dt_full,
					angka_max_win4d_disc, 
					angka_max_win3d_disc,angka_max_win3dd_disc, 
					angka_max_win2d_disc,angka_max_win2dd_disc,angka_max_win2dt_disc, 
					angka_max_win4d_bb, 
					angka_max_win3d_bb,angka_max_win3dd_bb, 
					angka_max_win2d_bb,angka_max_win2dd_bb,angka_max_win2dt_bb,
					angka_max_win4d_bb_kena,
					angka_max_win3d_bb_kena,angka_max_win3dd_bb_kena,
					angka_max_win2d_bb_kena,angka_max_win2dd_bb_kena,angka_max_win2dt_bb_kena,
                	create_by, create_at
				) 
              	VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
					$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
					$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, 
					$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, 
					$41, $42
				)`

	_, err := c.db.Exec(ctx, query,
		compconftoto.IDcompconftoto,
		compconftoto.IDcompany,
		compconftoto.AngkaMaxMinbasket,
		compconftoto.AngkaMaxMinbet,
		compconftoto.AngkaMaxMaxbet4d,
		compconftoto.AngkaMaxMaxbet3d,
		compconftoto.AngkaMaxMaxbet3dd,
		compconftoto.AngkaMaxMaxbet2d,
		compconftoto.AngkaMaxMaxbet2dd,
		compconftoto.AngkaMaxMaxbet2dt,
		compconftoto.AngkaMaxMaxbet4dBbdisc,
		compconftoto.AngkaMaxMaxbet3dBbdisc,
		compconftoto.AngkaMaxMaxbet3ddBbdisc,
		compconftoto.AngkaMaxMaxbet2dBbdisc,
		compconftoto.AngkaMaxMaxbet2ddBbdisc,
		compconftoto.AngkaMaxMaxbet2dtBbdisc,
		compconftoto.AngkaMaxWin4dFull,
		compconftoto.AngkaMaxWin3dFull,
		compconftoto.AngkaMaxWin3ddFull,
		compconftoto.AngkaMaxWin2dFull,
		compconftoto.AngkaMaxWin2ddFull,
		compconftoto.AngkaMaxWin2dtFull,
		compconftoto.AngkaMaxWin4dDisc,
		compconftoto.AngkaMaxWin3dDisc,
		compconftoto.AngkaMaxWin3ddDisc,
		compconftoto.AngkaMaxWin2dDisc,
		compconftoto.AngkaMaxWin2ddDisc,
		compconftoto.AngkaMaxWin2dtDisc,
		compconftoto.AngkaMaxWin4dBb,
		compconftoto.AngkaMaxWin3dBb,
		compconftoto.AngkaMaxWin3ddBb,
		compconftoto.AngkaMaxWin2dBb,
		compconftoto.AngkaMaxWin2ddBb,
		compconftoto.AngkaMaxWin2dtBb,
		compconftoto.AngkaMaxWin4dBbKena,
		compconftoto.AngkaMaxWin3dBbKena,
		compconftoto.AngkaMaxWin3ddBbKena,
		compconftoto.AngkaMaxWin2dBbKena,
		compconftoto.AngkaMaxWin2ddBbKena,
		compconftoto.AngkaMaxWin2dtBbKena,
		compconftoto.CreateBy,
		compconftoto.CreateAt,
	)
	return err
}

func (c companyconftotoRepository) UpdateAngka(ctx context.Context, compconftoto *domain.Companyconftoto) error {
	var query string
	var args []any

	query = `UPDATE ` + config.DB_tbl_companyconftoto + ` SET 
					angka_max_minbasket = $1, 
					angka_max_minbet = $2, 
					angka_max_maxbet_4d = $3, 
					angka_max_maxbet_3d = $4, 
					angka_max_maxbet_3dd = $5, 
					angka_max_maxbet_2d = $6, 
					angka_max_maxbet_2dd = $7, 
					angka_max_maxbet_2dt = $8, 
					angka_max_maxbet_4d_bbdisc = $9, 
					angka_max_maxbet_3d_bbdisc = $10, 
					angka_max_maxbet_3dd_bbdisc = $11,  
					angka_max_maxbet_2d_bbdisc = $12, 
					angka_max_maxbet_2dd_bbdisc = $13, 
					angka_max_maxbet_2dt_bbdisc = $14, 
					angka_max_win4d_full = $15, 
					angka_max_win3d_full = $16, 
					angka_max_win3dd_full = $17, 
					angka_max_win2d_full = $18,
					angka_max_win2dd_full = $19,
					angka_max_win2dt_full = $20,
					angka_max_win4d_disc = $21, 
					angka_max_win3d_disc = $22,
					angka_max_win3dd_disc = $23, 
					angka_max_win2d_disc = $24,
					angka_max_win2dd_disc = $25,
					angka_max_win2dt_disc = $26, 
					angka_max_win4d_bb = $27, 
					angka_max_win3d_bb = $28,
					angka_max_win3dd_bb = $29, 
					angka_max_win2d_bb = $30,
					angka_max_win2dd_bb = $31,
					angka_max_win2dt_bb = $32,
					angka_max_win4d_bb_kena = $33,
					angka_max_win3d_bb_kena = $34,
					angka_max_win3dd_bb_kena = $35,
					angka_max_win2d_bb_kena = $36,
					angka_max_win2dd_bb_kena = $37, 
					angka_max_win2dt_bb_kena = $38,
                    update_by = $39,
					update_at = $40
                  WHERE idcompconftoto = $41 AND idcompany = $42`
	args = []any{
		compconftoto.AngkaMaxMinbasket,
		compconftoto.AngkaMaxMinbet,
		compconftoto.AngkaMaxMaxbet4d,
		compconftoto.AngkaMaxMaxbet3d,
		compconftoto.AngkaMaxMaxbet3dd,
		compconftoto.AngkaMaxMaxbet2d,
		compconftoto.AngkaMaxMaxbet2dd,
		compconftoto.AngkaMaxMaxbet2dt,
		compconftoto.AngkaMaxMaxbet4dBbdisc,
		compconftoto.AngkaMaxMaxbet3dBbdisc,
		compconftoto.AngkaMaxMaxbet3ddBbdisc,
		compconftoto.AngkaMaxMaxbet2dBbdisc,
		compconftoto.AngkaMaxMaxbet2ddBbdisc,
		compconftoto.AngkaMaxMaxbet2dtBbdisc,
		compconftoto.AngkaMaxWin4dFull,
		compconftoto.AngkaMaxWin3dFull,
		compconftoto.AngkaMaxWin3ddFull,
		compconftoto.AngkaMaxWin2dFull,
		compconftoto.AngkaMaxWin2ddFull,
		compconftoto.AngkaMaxWin2dtFull,
		compconftoto.AngkaMaxWin4dDisc,
		compconftoto.AngkaMaxWin3dDisc,
		compconftoto.AngkaMaxWin3ddDisc,
		compconftoto.AngkaMaxWin2dDisc,
		compconftoto.AngkaMaxWin2ddDisc,
		compconftoto.AngkaMaxWin2dtDisc,
		compconftoto.AngkaMaxWin4dBb,
		compconftoto.AngkaMaxWin3dBb,
		compconftoto.AngkaMaxWin3ddBb,
		compconftoto.AngkaMaxWin2dBb,
		compconftoto.AngkaMaxWin2ddBb,
		compconftoto.AngkaMaxWin2dtBb,
		compconftoto.AngkaMaxWin4dBbKena,
		compconftoto.AngkaMaxWin3dBbKena,
		compconftoto.AngkaMaxWin3ddBbKena,
		compconftoto.AngkaMaxWin2dBbKena,
		compconftoto.AngkaMaxWin2ddBbKena,
		compconftoto.AngkaMaxWin2dtBbKena,
		compconftoto.UpdateBy,
		compconftoto.UpdateAt,
		compconftoto.IDcompconftoto,
		compconftoto.IDcompany,
	}

	res, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
