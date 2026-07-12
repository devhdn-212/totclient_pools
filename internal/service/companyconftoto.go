package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/repository"
	"github.com/devhdn-212/totagen_api/internal/util"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisCompanyconftotoKey = "master:companyconftoto:all"
)

type companyconftotoService struct {
	db   *pgxpool.Pool
	repo domain.CompanyconftotoRepository
}

func NewCompanyconftotoService(db *pgxpool.Pool, repo domain.CompanyconftotoRepository) domain.CompanyconftotoService {
	return &companyconftotoService{
		db:   db,
		repo: repo,
	}
}
func (c companyconftotoService) All(ctx context.Context, idcomp string) ([]dto.CompanyconftotoData, error) {
	cached, found, err := connection.GetRedis(RedisCompanyconftotoKey + strings.ToLower(idcomp))
	if err != nil {
		return nil, err
	}
	var data []dto.CompanyconftotoData
	if found {

		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Company Conf Toto")
			return data, nil
		}
	}

	record, err := c.repo.FindAll(ctx, idcomp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range record {
		var createdAt string

		if v.CreateAt.Valid {
			if v.CreateBy != "" {
				createdAt = v.CreateBy + ", " + v.CreateAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
			} else {
				createdAt = ""
			}
		}

		data = append(data, dto.CompanyconftotoData{
			IDcompconftoto:          v.IDcompconftoto,
			IDcompany:               v.IDcompany,
			AngkaMaxMinbasket:       v.AngkaMaxMinbasket,
			AngkaMaxMinbet:          v.AngkaMaxMinbet,
			AngkaMaxMaxbet4d:        v.AngkaMaxMaxbet4d,
			AngkaMaxMaxbet3d:        v.AngkaMaxMaxbet3d,
			AngkaMaxMaxbet3dd:       v.AngkaMaxMaxbet3dd,
			AngkaMaxMaxbet2d:        v.AngkaMaxMaxbet2d,
			AngkaMaxMaxbet2dd:       v.AngkaMaxMaxbet2dd,
			AngkaMaxMaxbet2dt:       v.AngkaMaxMaxbet2dt,
			AngkaMaxMaxbet4dBbdisc:  v.AngkaMaxMaxbet4dBbdisc,
			AngkaMaxMaxbet3dBbdisc:  v.AngkaMaxMaxbet3dBbdisc,
			AngkaMaxMaxbet3ddBbdisc: v.AngkaMaxMaxbet3ddBbdisc,
			AngkaMaxMaxbet2dBbdisc:  v.AngkaMaxMaxbet2dBbdisc,
			AngkaMaxMaxbet2ddBbdisc: v.AngkaMaxMaxbet2ddBbdisc,
			AngkaMaxMaxbet2dtBbdisc: v.AngkaMaxMaxbet2dtBbdisc,
			AngkaMaxWin4dFull:       v.AngkaMaxWin4dFull,
			AngkaMaxWin3dFull:       v.AngkaMaxWin3dFull,
			AngkaMaxWin3ddFull:      v.AngkaMaxWin3ddFull,
			AngkaMaxWin2dFull:       v.AngkaMaxWin2dFull,
			AngkaMaxWin2ddFull:      v.AngkaMaxWin2ddFull,
			AngkaMaxWin2dtFull:      v.AngkaMaxWin2dtFull,
			AngkaMaxWin4dDisc:       v.AngkaMaxWin4dDisc,
			AngkaMaxWin3dDisc:       v.AngkaMaxWin3dDisc,
			AngkaMaxWin3ddDisc:      v.AngkaMaxWin3ddDisc,
			AngkaMaxWin2dDisc:       v.AngkaMaxWin2dDisc,
			AngkaMaxWin2ddDisc:      v.AngkaMaxWin2ddDisc,
			AngkaMaxWin2dtDisc:      v.AngkaMaxWin2dtDisc,
			AngkaMaxWin4dBb:         v.AngkaMaxWin4dBb,
			AngkaMaxWin3dBb:         v.AngkaMaxWin3dBb,
			AngkaMaxWin3ddBb:        v.AngkaMaxWin3ddBb,
			AngkaMaxWin2dBb:         v.AngkaMaxWin2dBb,
			AngkaMaxWin2ddBb:        v.AngkaMaxWin2ddBb,
			AngkaMaxWin2dtBb:        v.AngkaMaxWin2dtBb,
			AngkaMaxWin4dBbKena:     v.AngkaMaxWin4dBbKena,
			AngkaMaxWin3dBbKena:     v.AngkaMaxWin3dBbKena,
			AngkaMaxWin3ddBbKena:    v.AngkaMaxWin3ddBbKena,
			AngkaMaxWin2dBbKena:     v.AngkaMaxWin2dBbKena,
			AngkaMaxWin2ddBbKena:    v.AngkaMaxWin2ddBbKena,
			AngkaMaxWin2dtBbKena:    v.AngkaMaxWin2dtBbKena,
			CbebasMaxMinbet:         v.CbebasMaxMinbet,
			CbebasMaxMaxbet:         v.CbebasMaxMaxbet,
			CbebasMaxWin:            v.CbebasMaxWin,
			CmacauMaxMinbet:         v.CmacauMaxMinbet,
			CmacauMaxMaxbet:         v.CmacauMaxMaxbet,
			CmacauMaxWin2:           v.CmacauMaxWin2,
			CmacauMaxWin3:           v.CmacauMaxWin3,
			CmacauMaxWin4:           v.CmacauMaxWin4,
			CnagaMaxMinbet:          v.CnagaMaxMinbet,
			CnagaMaxMaxbet:          v.CnagaMaxMaxbet,
			CnagaMaxWin3:            v.CnagaMaxWin3,
			CnagaMaxWin4:            v.CnagaMaxWin4,
			CjituMaxMinbet:          v.CjituMaxMinbet,
			CjituMaxMaxbet:          v.CjituMaxMaxbet,
			CjituMaxWinas:           v.CjituMaxWinas,
			CjituMaxWinkop:          v.CjituMaxWinkop,
			CjituMaxWinkepala:       v.CjituMaxWinkepala,
			CjituMaxWinekor:         v.CjituMaxWinekor,
			Umum50MaxMinbet:         v.Umum50MaxMinbet,
			Umum50MaxMaxbet:         v.Umum50MaxMaxbet,
			Special50MaxMinbet:      v.Special50MaxMinbet,
			Special50MaxMaxbet:      v.Special50MaxMaxbet,
			Kombinasi50MaxMinbet:    v.Kombinasi50MaxMinbet,
			Kombinasi50MaxMaxbet:    v.Kombinasi50MaxMaxbet,
			MacauMaxMinbet:          v.MacauMaxMinbet,
			MacauMaxMaxbet:          v.MacauMaxMaxbet,
			MacauMaxWin:             v.MacauMaxWin,
			DasarMaxMinbet:          v.DasarMaxMinbet,
			DasarMaxMaxbet:          v.DasarMaxMaxbet,
			ShioMaxMinbet:           v.ShioMaxMinbet,
			ShioMaxMaxbet:           v.ShioMaxMaxbet,
			ShioMaxWin:              v.ShioMaxWin,
			ShioParent:              v.ShioParent,
			CreateBy:                createdAt,
		})
	}

	go connection.SetRedis(RedisCompanyconftotoKey+strings.ToLower(idcomp), data, 60*time.Minute)
	connection.Log.Info("Returning data Database - Company Conf TOTO")
	return data, nil
}

func (c companyconftotoService) Save(ctx context.Context, req dto.CompanyconftotoSave, client_admin string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewCompanyconftotoRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.IDcompany)
	if err != nil {
		return err
	}

	now := util.GetNowJakarta()

	if req.Type == "Edit" {
		if flag.IDcompconftoto == "" {
			return fmt.Errorf("Company conf toto %w", util.ErrNotFound)
		}
		flag.IDcompconftoto = req.IDcompconftoto
		flag.IDcompany = req.IDcompany
		flag.AngkaMaxMinbasket = req.AngkaMaxMinbasket
		flag.AngkaMaxMinbet = req.AngkaMaxMinbet
		flag.AngkaMaxMaxbet4d = req.AngkaMaxMaxbet4d
		flag.AngkaMaxMaxbet3d = req.AngkaMaxMaxbet3d
		flag.AngkaMaxMaxbet3dd = req.AngkaMaxMaxbet3dd
		flag.AngkaMaxMaxbet2d = req.AngkaMaxMaxbet2d
		flag.AngkaMaxMaxbet2dd = req.AngkaMaxMaxbet2dd
		flag.AngkaMaxMaxbet2dt = req.AngkaMaxMaxbet2dt
		flag.AngkaMaxMaxbet4dBbdisc = req.AngkaMaxMaxbet4dBbdisc
		flag.AngkaMaxMaxbet3dBbdisc = req.AngkaMaxMaxbet3dBbdisc
		flag.AngkaMaxMaxbet3ddBbdisc = req.AngkaMaxMaxbet3ddBbdisc
		flag.AngkaMaxMaxbet2dBbdisc = req.AngkaMaxMaxbet2dBbdisc
		flag.AngkaMaxMaxbet2ddBbdisc = req.AngkaMaxMaxbet2ddBbdisc
		flag.AngkaMaxMaxbet2dtBbdisc = req.AngkaMaxMaxbet2dtBbdisc
		flag.AngkaMaxWin4dFull = req.AngkaMaxWin4dFull
		flag.AngkaMaxWin3dFull = req.AngkaMaxWin3dFull
		flag.AngkaMaxWin3ddFull = req.AngkaMaxWin3ddFull
		flag.AngkaMaxWin2dFull = req.AngkaMaxWin2dFull
		flag.AngkaMaxWin2ddFull = req.AngkaMaxWin2ddFull
		flag.AngkaMaxWin2dtFull = req.AngkaMaxWin2dtFull
		flag.AngkaMaxWin4dDisc = req.AngkaMaxWin4dDisc
		flag.AngkaMaxWin3dDisc = req.AngkaMaxWin3dDisc
		flag.AngkaMaxWin3ddDisc = req.AngkaMaxWin3ddDisc
		flag.AngkaMaxWin2dDisc = req.AngkaMaxWin2dDisc
		flag.AngkaMaxWin2ddDisc = req.AngkaMaxWin2ddDisc
		flag.AngkaMaxWin2dtDisc = req.AngkaMaxWin2dtDisc
		flag.AngkaMaxWin4dBb = req.AngkaMaxWin4dBb
		flag.AngkaMaxWin3dBb = req.AngkaMaxWin3dBb
		flag.AngkaMaxWin3ddBb = req.AngkaMaxWin3ddBb
		flag.AngkaMaxWin2dBb = req.AngkaMaxWin2dBb
		flag.AngkaMaxWin2ddBb = req.AngkaMaxWin2ddBb
		flag.AngkaMaxWin2dtBb = req.AngkaMaxWin2dtBb
		flag.AngkaMaxWin4dBbKena = req.AngkaMaxWin4dBbKena
		flag.AngkaMaxWin3dBbKena = req.AngkaMaxWin3dBbKena
		flag.AngkaMaxWin3ddBbKena = req.AngkaMaxWin3ddBbKena
		flag.AngkaMaxWin2dBbKena = req.AngkaMaxWin2dBbKena
		flag.AngkaMaxWin2ddBbKena = req.AngkaMaxWin2ddBbKena
		flag.AngkaMaxWin2dtBbKena = req.AngkaMaxWin2dtBbKena
		flag.CbebasMaxMinbet = req.CbebasMaxMinbet
		flag.CbebasMaxMaxbet = req.CbebasMaxMaxbet
		flag.CbebasMaxWin = req.CbebasMaxWin
		flag.CmacauMaxMinbet = req.CmacauMaxMinbet
		flag.CmacauMaxMaxbet = req.CmacauMaxMaxbet
		flag.CmacauMaxWin2 = req.CmacauMaxWin2
		flag.CmacauMaxWin3 = req.CmacauMaxWin3
		flag.CmacauMaxWin4 = req.CmacauMaxWin4
		flag.CnagaMaxMinbet = req.CnagaMaxMinbet
		flag.CnagaMaxMaxbet = req.CnagaMaxMaxbet
		flag.CnagaMaxWin3 = req.CnagaMaxWin3
		flag.CnagaMaxWin4 = req.CnagaMaxWin4
		flag.CjituMaxMinbet = req.CjituMaxMinbet
		flag.CjituMaxMaxbet = req.CjituMaxMaxbet
		flag.CjituMaxWinas = req.CjituMaxWinas
		flag.CjituMaxWinkop = req.CjituMaxWinkop
		flag.CjituMaxWinkepala = req.CjituMaxWinkepala
		flag.CjituMaxWinekor = req.CjituMaxWinekor
		flag.Umum50MaxMinbet = req.Umum50MaxMinbet
		flag.Umum50MaxMaxbet = req.Umum50MaxMaxbet
		flag.Special50MaxMinbet = req.Special50MaxMinbet
		flag.Special50MaxMaxbet = req.Special50MaxMaxbet
		flag.Kombinasi50MaxMinbet = req.Kombinasi50MaxMinbet
		flag.Kombinasi50MaxMaxbet = req.Kombinasi50MaxMaxbet
		flag.MacauMaxMinbet = req.MacauMaxMinbet
		flag.MacauMaxMaxbet = req.MacauMaxMaxbet
		flag.MacauMaxWin = req.MacauMaxWin
		flag.DasarMaxMinbet = req.DasarMaxMinbet
		flag.DasarMaxMaxbet = req.DasarMaxMaxbet
		flag.ShioMaxMinbet = req.ShioMaxMinbet
		flag.ShioMaxMaxbet = req.ShioMaxMaxbet
		flag.ShioMaxWin = req.ShioMaxWin
		flag.ShioParent = req.ShioParent
		flag.UpdateBy = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		if err = txRepo.Update(ctx, &flag); err != nil {
			return err
		}

	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisCompanyconftotoKey + strings.ToLower(req.IDcompany))
	return nil
}
