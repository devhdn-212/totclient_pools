package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/repository"
	"github.com/devhdn-212/totclient_api/internal/util"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisAdminAllKey = "agen:admin:all"
)

type adminService struct {
	db   *pgxpool.Pool
	repo domain.AdminsRepository
}

func NewAdminService(db *pgxpool.Pool, repo domain.AdminsRepository) domain.AdminService {
	return &adminService{
		db:   db,
		repo: repo,
	}
}
func (a adminService) All(ctx context.Context, idcomp string) ([]dto.AdminData, error) {
	cached, found, err := connection.GetRedis(RedisAdminAllKey + ":" + strings.ToLower(idcomp))
	if err != nil {
		return nil, err
	}
	var data []dto.AdminData
	if found {
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Admin")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	admins, err := a.repo.FindAll(ctx, idcomp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range admins {
		var lastlogin, createdAt, updatedAt string

		if v.Lastlogin.Valid {
			lastlogin = v.Lastlogin.Time.In(util.LocJakarta).Format("2006-01-02 15:04:05")
		}
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

		data = append(data, dto.AdminData{
			ID:        v.ID,
			Idrule:    v.IDClientrule,
			Username:  v.Username,
			Name:      v.Name,
			Status:    v.Status,
			Lastlogin: lastlogin,
			Ipaddress: v.Ipaddress,
			Created:   createdAt,
			Update:    updatedAt,
		})
	}

	go connection.SetRedis(RedisAdminAllKey+":"+strings.ToLower(idcomp), data, 60*time.Minute)
	connection.Log.Info("Returning data Database - Admin")
	return data, nil
}

func (a adminService) Save(ctx context.Context, req dto.AdminSave, client_admin, idcomp string) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewAdminRepository(txExec)

	flag, err := txRepo.FindByUsernameComp(ctx, req.Username, strings.ToUpper(idcomp))
	if err != nil {
		return err
	}

	// Set lokasi ke Asia/Jakarta
	now := util.GetNowJakarta()

	haspass, _ := util.HashPassword(req.Pass)

	if req.Type == "New" {
		if flag.Username != "" {
			return util.ErrDuplicate
		}

		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		date := time.Now().Format("0601")
		idcompadmin := fmt.Sprintf("%s-%s-admin-%s", strings.ToLower(idcomp), date, raw)
		username := strings.ToLower(idcomp) + req.Username

		admin := domain.Admin{
			ID:           idcompadmin,
			IDClientrule: req.IDrule,
			IDCompany:    strings.ToUpper(idcomp),
			Username:     username,
			Pass:         haspass,
			Name:         req.Name,
			Status:       req.Status,
			Created:      client_admin,
			CreatedAt:    sql.NullTime{Valid: true, Time: now},
		}

		err = txRepo.Save(ctx, &admin)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return util.ErrDuplicate
			}
			return err
		}
	} else {
		if flag.ID == "" {
			return errors.New("Username not found")
		}

		flag.ID = req.ID
		flag.IDCompany = strings.ToUpper(idcomp)
		flag.IDClientrule = req.IDrule
		flag.Name = req.Name
		flag.Status = req.Status
		flag.Update = client_admin
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		if req.Pass != "" {
			flag.Pass = haspass
			if err = txRepo.Update(ctx, &flag, true); err != nil {
				return err
			}
		} else {
			if err = txRepo.Update(ctx, &flag, false); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisAdminAllKey + ":" + strings.ToLower(idcomp))
	return nil
}
