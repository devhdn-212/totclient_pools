package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/devhdn-212/gofibermaster_api/domain"
	"github.com/devhdn-212/gofibermaster_api/dto"
	"github.com/devhdn-212/gofibermaster_api/internal/connection"
	"github.com/devhdn-212/gofibermaster_api/internal/repository"
	"github.com/devhdn-212/gofibermaster_api/internal/util"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	RedisAdminAllKey = "master:admin:all"
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
func (a adminService) All(ctx context.Context) ([]dto.AdminData, error) {
	cached, found, err := connection.GetRedis(RedisAdminAllKey)
	if err != nil {
		return nil, err
	}

	if found {
		var data []dto.AdminData
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			connection.Log.Info("Returning data from Redis - Admin")
			return data, nil
		}
		// kalau corrupt → lanjut ke DB
	}

	admins, err := a.repo.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var adminData []dto.AdminData
	for _, v := range admins {
		var joindate, lastlogin, createdAt, updatedAt string

		if v.Joindate.Valid {
			joindate = v.Joindate.Time.In(util.LocJakarta).Format("2006-01-02")
		}
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

		adminData = append(adminData, dto.AdminData{
			Username:  v.Username,
			Idadmin:   v.Idadmin,
			Name:      v.Name,
			Status:    v.Status,
			Lastlogin: lastlogin,
			Joindate:  joindate,
			Ipaddress: v.Ipaddress,
			Created:   createdAt,
			Update:    updatedAt,
		})
	}

	go connection.SetRedis(RedisAdminAllKey, adminData, 60*time.Minute)
	connection.Log.Info("Returning data Database - Admin")
	return adminData, nil
}

func (a adminService) Save(ctx context.Context, req dto.AdminSave, client_admin string) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewAdminRepository(txExec)

	flag, err := txRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	// Set lokasi ke Asia/Jakarta
	now := util.GetNowJakarta()

	haspass, _ := util.HashPassword(req.Pass)

	if req.Type == "New" {
		if flag.Username != "" {
			return errors.New("Username already exists")
		}

		admin := domain.Admin{
			Username: req.Username,
			Pass:     haspass,
			Idadmin:  req.Idadmin,
			Name:     req.Name,
			Status:   req.Status,
			// Input waktu dengan timezone JKT
			Lastlogin: sql.NullTime{Valid: true, Time: now},
			Joindate:  sql.NullTime{Valid: true, Time: now},
			Ipaddress: req.Ipaddress,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: now},
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
		if flag.Username == "" {
			return errors.New("Username not found")
		}

		flag.Name = req.Name
		flag.Status = req.Status
		flag.Ipaddress = req.Ipaddress
		flag.Update = client_admin
		// Update waktu dengan timezone JKT
		flag.UpdateAt = sql.NullTime{Valid: true, Time: now}

		if req.Pass != "" {
			flag.Pass = haspass
			if err = txRepo.Update(ctx, &flag); err != nil {
				return err
			}
		} else {
			if err = txRepo.UpdateLogin(ctx, &flag); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	go connection.DeleteRedis(RedisAdminAllKey)
	return nil
}
