package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/util"

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
