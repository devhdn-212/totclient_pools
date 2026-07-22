package repository

import (
	"context"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/util"
)

type trxkeluarandetailrejectRepository struct {
	db DBExecutor
}

func NewTrxkeluarandetailrejectRepository(db DBExecutor) domain.TrxkeluarandetailrejectRepository {
	return &trxkeluarandetailrejectRepository{
		db: db,
	}
}

func (a trxkeluarandetailrejectRepository) Save(ctx context.Context, reject *domain.Trxkeluarandetailreject, idcomp string) error {
	t := util.Get_mapping_totodb(idcomp)
	query := `INSERT INTO ` + t.Schema + `.` + t.KeluarantogelDetailReject + `
                (idtrxkeluarandetailreject, idtrxkeluaran, idcompany,
				datetimedetail, ipaddress,
				username, typegame, posisitogel, nomortogel,
				bet, playerinvoice, reason, sisalimit,
				browsertogel, devicetogel,
				create_by, create_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := a.db.Exec(ctx, query,
		reject.ID,
		reject.IDtrxkeluaran,
		reject.IDcomp,
		reject.Datekeluarandetail,
		reject.Ipaddress,
		reject.Username,
		reject.Typegame,
		reject.Posisitogel,
		reject.Nomortogel,
		reject.Bet,
		reject.Playerinvoice,
		reject.Reason,
		reject.Sisalimit,
		reject.Browsertogel,
		reject.Devicetogel,
		reject.Created,
		reject.CreatedAt,
	)
	return err
}
