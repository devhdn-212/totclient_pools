package repository

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/jackc/pgx/v5"
)

// trxkeluaranmemberColumns lists exactly the columns
// domain.Trxkeluaranmember has fields for — see trxkeluarandetailColumns for
// why SELECT * isn't used here.
const trxkeluaranmemberColumns = `idkeluaranmember, idtrxkeluaran, idcompany,
	username, totalbet, totalbayar, totaldiscount, totalkei, totalwin,
	totalpair, betround, playerinvoice, status,
	createkeluaranmember, createdatekeluaranmember,
	COALESCE(updatekeluaranmember, '') AS updatekeluaranmember, updatedatekeluaranmember`

type trxkeluaranmemberRepository struct {
	db DBExecutor
}

func NewTrxkeluaranmemberRepository(db DBExecutor) domain.TrxkeluaranmemberRepository {
	return &trxkeluaranmemberRepository{
		db: db,
	}
}
func (a trxkeluaranmemberRepository) FindAll(ctx context.Context, idcomp string, idtrx int) ([]domain.Trxkeluaranmember, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluaranmemberColumns + ` FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			WHERE idtrxkeluaran = $1
			ORDER BY createdatekeluaranmember DESC`

	rows, err := a.db.Query(ctx, query, idtrx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Trxkeluaranmember])
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (a trxkeluaranmemberRepository) FindAllByUsername(ctx context.Context, idcomp string, idtrx int, username string) ([]domain.Trxkeluaranmember, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluaranmemberColumns + ` FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			WHERE idtrxkeluaran = $1 AND username = $2
			ORDER BY createdatekeluaranmember DESC
			LIMIT 100`

	rows, err := a.db.Query(ctx, query, idtrx, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Trxkeluaranmember])
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (a trxkeluaranmemberRepository) FindByID(ctx context.Context, idcomp, username string, idtrx int) (domain.Trxkeluaranmember, error) {
	t := util.Get_mapping_totodb(idcomp)
	query := `SELECT ` + trxkeluaranmemberColumns + `
			FROM ` + t.Schema + `.` + t.KeluarantogelMember + `
			WHERE username = $1 AND idtrxkeluaran = $2
			LIMIT 1`

	rows, err := a.db.Query(ctx, query, idtrx, username, idtrx)
	if err != nil {
		return domain.Trxkeluaranmember{}, err
	}
	defer rows.Close()

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[domain.Trxkeluaranmember])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trxkeluaranmember{}, nil
		}
		return domain.Trxkeluaranmember{}, err
	}
	return data, nil
}
func (a trxkeluaranmemberRepository) Save(ctx context.Context, trxkeluaranmember *domain.Trxkeluaranmember, idcomp string) error {
	t := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + t.Schema + `.` + t.KeluarantogelMember + `
                (idkeluaranmember, idtrxkeluaran, idcompany, 
				username, 
				totalbet, totalbayar, totaldiscount, totalkei, totalwin, totalpair, betround,
				playerinvoice, status, 
				createkeluaranmember, createdatekeluaranmember) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13,$14,$15)`

	_, err := a.db.Exec(ctx, query,
		trxkeluaranmember.ID,
		trxkeluaranmember.IDtrxkeluaran,
		trxkeluaranmember.IDcomp,
		trxkeluaranmember.Username,
		trxkeluaranmember.Totalbet,
		trxkeluaranmember.Totalbayar,
		trxkeluaranmember.Totaldiscount,
		trxkeluaranmember.Totalkei,
		trxkeluaranmember.Totalwin,
		trxkeluaranmember.Totalpair,
		trxkeluaranmember.Betround,
		trxkeluaranmember.Playerinvoice,
		trxkeluaranmember.Status,
		trxkeluaranmember.Created,
		trxkeluaranmember.CreatedAt,
	)
	return err
}

func (a trxkeluaranmemberRepository) Update(ctx context.Context, trxkeluaranmember *domain.Trxkeluaranmember, idcomp string) error {
	var query string
	var args []any

	t := util.Get_mapping_totodb(idcomp)
	query = `UPDATE ` + t.Schema + `.` + t.KeluarantogelMember + ` SET
                    totalwin = $1, 
                    status = $2, 
                    updatekeluaranmember = $3, 
                    updatedatekeluaranmember = $4  
                  WHERE username = $5 AND playerinvoice =$6 AND idtrxkeluaran=$7 `
	args = []any{
		trxkeluaranmember.Totalwin,
		trxkeluaranmember.Status,
		trxkeluaranmember.Update,
		trxkeluaranmember.UpdateAt,
		trxkeluaranmember.Username,
		trxkeluaranmember.Playerinvoice,
		trxkeluaranmember.IDtrxkeluaran,
	}

	res, err := a.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
