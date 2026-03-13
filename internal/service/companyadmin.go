package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/dto"
	"github.com/devhdn-212/gofibergoqu_master/internal/connection"
	"github.com/devhdn-212/gofibergoqu_master/internal/repository"
	"github.com/devhdn-212/gofibergoqu_master/internal/util"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	RedisCompanyadminKey = "master:companyadmin:"
)

type companyadminService struct {
	db   *sql.DB
	repo domain.CompanyadminRepository
}

func NewCompanyadminService(db *sql.DB, repo domain.CompanyadminRepository) domain.CompanyadminService {
	return &companyadminService{
		db:   db,
		repo: repo,
	}
}
func (c companyadminService) All(ctx context.Context, idcompany string) ([]dto.CompanyadminData, error) {
	cached, found, err := connection.GetRedis(RedisCompanyadminKey + strings.ToLower(idcompany))
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.CompanyadminData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Company Admin")
			return data, nil
		}
	}

	rescompadmin, err := c.repo.FindAll(ctx, idcompany)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var compadminData []dto.CompanyadminData
	for _, v := range rescompadmin {
		var lastlogin, createdAt, updatedAt string
		if v.Lastlogin.Valid {
			lastlogin = v.Lastlogin.Time.In(loc).Format("2006-01-02 15:04:05")
		}
		if v.CreatedAt.Valid {
			createdAt = v.Created + ", " + v.CreatedAt.Time.In(loc).Format("2006-01-02 15:04:05")
		}
		if v.UpdateAt.Valid {
			if v.Update != "" {
				updatedAt = v.Update + ", " + v.UpdateAt.Time.In(loc).Format("2006-01-02 15:04:05")
			} else {
				updatedAt = ""
			}
		}

		compadminData = append(compadminData, dto.CompanyadminData{
			ID:        v.ID,
			IDcomp:    v.IDCompany,
			Rule:      v.IDClientrule,
			Username:  v.Username,
			Name:      v.Name,
			Lastlogin: lastlogin,
			Ipaddress: v.Ipaddress,
			Status:    v.Status,
			Created:   createdAt,
			Update:    updatedAt,
		})
	}

	go connection.SetRedis(RedisCompanyadminKey+strings.ToLower(idcompany), compadminData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Company Admin")
	return compadminData, nil
}

func (c companyadminService) Save(ctx context.Context, req dto.CompanyadminSave, client string) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewCompanyadminRepository(txExec)
	flag, err := txRepo.FindByID(ctx, req.IDcompany, req.Username)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	haspass, _ := util.HashPassword(req.Pass)
	if req.Type == "New" {
		if flag.ID != "" {
			return errors.New("Duplicate Entry")
		}
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		date := time.Now().Format("0601")
		idcompadmin := fmt.Sprintf("%s-%s-admin-%s", strings.ToLower(req.IDcompany), date, raw)
		comp := domain.Companyadmin{
			ID:           idcompadmin,
			IDCompany:    req.IDcompany,
			IDClientrule: req.IDrule,
			Username:     req.Username,
			Pass:         haspass,
			Name:         req.Name,
			Status:       req.Status,
			Created:      client,
			CreatedAt:    sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &comp)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	} else {
		if flag.ID == "" {
			return errors.New("Username not found")
		}
		if req.Pass != "" {
			flag.ID = req.ID
			flag.IDClientrule = req.IDrule
			flag.Pass = haspass
			flag.Name = req.Name
			flag.Status = req.Status
			flag.Update = client
			flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

			if err = c.repo.Update(ctx, &flag, true); err != nil {
				fmt.Println(err)
				return err
			}
			if err := tx.Commit(); err != nil {
				return err
			}
		} else {
			flag.ID = req.ID
			flag.IDClientrule = req.IDrule
			flag.Name = req.Name
			flag.Status = req.Status
			flag.Update = client
			flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

			if err = c.repo.Update(ctx, &flag, false); err != nil {
				fmt.Println(err)
				return err
			}
			if err := tx.Commit(); err != nil {
				return err
			}
		}
	}

	go connection.DeleteRedis(RedisCompanyadminKey + strings.ToLower(req.IDcompany))
	return nil
}
