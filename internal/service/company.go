package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/dto"
	"github.com/devhdn-212/totmaster_api/internal/connection"
	"github.com/devhdn-212/totmaster_api/internal/repository"
	"github.com/devhdn-212/totmaster_api/internal/util"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgconn"
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

func (c companyService) Save(ctx context.Context, req dto.CompanySave, client_admin string) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewCompanyRepository(txExec)
	txRepoConfToto := repository.NewCompanyconftotoRepository(txExec)

	flag, err := txRepo.FindByID(ctx, req.ID)
	flag_conf, err_conf := txRepoConfToto.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if err_conf != nil {
		return err_conf
	}

	now := util.GetNowJakarta()

	if req.Type == "New" {
		if flag.ID != "" {
			return util.ErrDuplicate
		}
		comp := domain.Company{
			ID:          req.ID,
			IDcurrdef:   req.IDcurr,
			IDgroupcomp: req.IDgroupcomp,
			Name:        req.Name,
			TelegramID:  req.TelegramID,
			URLapitoto:  req.URLapitoto,
			URLapislot:  req.URLapislot,
			Status:      req.Status,
			Activetoto:  req.Activetoto,
			Activeslot:  req.Activeslot,
			Created:     client_admin,
			CreatedAt:   sql.NullTime{Valid: true, Time: now},
		}
		err = txRepo.Save(ctx, &comp)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
		if flag.Activetoto == "Y" {
			if flag_conf.IDcompconftoto == "" { // TOTO ACTIVE
				fmt.Println("==CREATE CONF TOTO==")
				raw := strings.ReplaceAll(uuid.NewString(), "-", "")
				date := time.Now().Format("0601")
				idconf := fmt.Sprintf("%s-%s-conftoto-%s", strings.ToLower(req.ID), date, raw)
				compconf := domain.Companyconftoto{
					IDcompconftoto: idconf,
					IDcompany:      req.ID,
					CreateBy:       client_admin,
					CreateAt:       sql.NullTime{Valid: true, Time: now},
				}
				err = txRepoConfToto.Save(ctx, &compconf)
				if err != nil {
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) && pgErr.Code == "23505" {
						return util.ErrDuplicate
					}
					return err
				}
				fmt.Println("==END CREATE CONF TOTO==")
			}
		}
	} else {
		if flag.ID == "" {
			return fmt.Errorf("Company %w", util.ErrNotFound)
		}

		flag.IDcurrdef = req.IDcurr
		flag.IDgroupcomp = req.IDgroupcomp
		flag.Name = req.Name
		flag.TelegramID = req.TelegramID
		flag.URLapitoto = req.URLapitoto
		flag.URLapislot = req.URLapislot
		flag.Status = req.Status
		flag.Activetoto = req.Activetoto
		flag.Activeslot = req.Activeslot
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		if err = txRepo.Update(ctx, &flag); err != nil {
			return err
		}
		if flag.Activetoto == "Y" {
			if flag_conf.IDcompconftoto == "" { // TOTO ACTIVE
				fmt.Println("==CREATE CONF TOTO==")
				raw := strings.ReplaceAll(uuid.NewString(), "-", "")
				date := time.Now().Format("0601")
				idconf := fmt.Sprintf("%s-%s-conftoto-%s", strings.ToLower(req.ID), date, raw)
				compconf := domain.Companyconftoto{
					IDcompconftoto: idconf,
					IDcompany:      req.ID,
					CreateBy:       client_admin,
					CreateAt:       sql.NullTime{Valid: true, Time: now},
				}
				err = txRepoConfToto.Save(ctx, &compconf)
				if err != nil {
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) && pgErr.Code == "23505" {
						return util.ErrDuplicate
					}
					return err
				}
				fmt.Println("==END CREATE CONF TOTO==")
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisCompanyKey)
	return nil
}
