package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/util"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisCompanyKey = "master:company:all"
)

type companyService struct {
	db   *pgxpool.Pool
	repo domain.CompanyRepository
}

func NewCompanyService(db *pgxpool.Pool, repo domain.CompanyRepository) domain.CompanyService {
	return &companyService{
		db:   db,
		repo: repo,
	}
}
func (c companyService) All(ctx context.Context) ([]dto.CompanyData, error) {
	cached, found, err := connection.GetRedis(RedisCompanyKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CompanyData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Company")
			return data, nil
		}
	}

	curr, err := c.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var compData []dto.CompanyData
	for _, v := range curr {
		var createdAt, updatedAt, endjoin string
		if v.CreatedAt.Valid {
			createdAt = v.Created + ", " + v.CreatedAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
		}
		if v.UpdateAt.Valid {
			if v.Update != "" {
				updatedAt = v.Update + ", " + v.UpdateAt.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
			} else {
				updatedAt = ""
			}
		}
		if v.Endjoin.Valid {
			endjoin = v.Endjoin.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
		}

		compData = append(compData, dto.CompanyData{
			ID:          v.ID,
			IDgroupcomp: v.IDgroupcomp,
			Nmgroupcomp: v.Nmgroupcomp.String,
			IDcurrdef:   v.IDcurrdef,
			Name:        v.Name,
			TelegramID:  v.TelegramID,
			URLapitoto:  v.URLapitoto,
			URLapislot:  v.URLapislot,
			Endjoin:     endjoin,
			Activetoto:  v.Activetoto,
			Activeslot:  v.Activeslot,
			Status:      v.Status,
			Created:     createdAt,
			Update:      updatedAt,
		})
	}

	go connection.SetRedis(RedisCompanyKey, compData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Company")
	return compData, nil
}
