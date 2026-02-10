package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/repository"
	"gofibergocu/internal/util"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/lib/pq"
)

const (
	RedisAdminAllKey = "admin:all"
)

type adminService struct {
	db   *sql.DB
	repo domain.AdminsRepository
}

func NewAdminService(db *sql.DB, repo domain.AdminsRepository) domain.AdminService {
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
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var adminData []dto.AdminData
	for _, v := range admins {
		var joindate, lastlogin, createdAt, updatedAt string

		if v.Joindate.Valid {
			joindate = v.Joindate.Time.In(loc).Format("2006-01-02")
		}
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
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewAdminRepository(txExec)
	flag, err := txRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	fmt.Println("Go time:", now)
	fmt.Println("Location:", now.Location())
	haspass, _ := util.HashPassword(req.Pass)
	if req.Type == "New" {
		if flag.Username != "" {
			return errors.New("Username already exists")
		}
		admin := domain.Admin{
			Username:  req.Username,
			Pass:      haspass,
			Idadmin:   req.Idadmin,
			Name:      req.Name,
			Status:    req.Status,
			Lastlogin: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
			Joindate:  sql.NullTime{Valid: true, Time: time.Now().In(loc)},
			Ipaddress: req.Ipaddress,
			Created:   client_admin,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now().In(loc)},
		}
		err = txRepo.Save(ctx, &admin)
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
		if flag.Username == "" {
			return errors.New("Username not found")
		}
		if req.Pass != "" {
			flag.Username = req.Username
			flag.Pass = haspass
			flag.Idadmin = req.Idadmin
			flag.Name = req.Name
			flag.Status = req.Status
			flag.Ipaddress = req.Status
			flag.Update = client_admin
			flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

			if err = a.repo.Update(ctx, &flag); err != nil {
				fmt.Println(err)
				return err
			}
			if err := tx.Commit(); err != nil {
				return err
			}
		} else {
			flag.Username = req.Username
			flag.Idadmin = req.Idadmin
			flag.Name = req.Name
			flag.Status = req.Status
			flag.Ipaddress = req.Status
			flag.Update = req.Username
			flag.UpdateAt = sql.NullTime{Valid: true, Time: time.Now().In(loc)}

			if err = a.repo.Update(ctx, &flag); err != nil {
				return err
			}
			if err := tx.Commit(); err != nil {
				return err
			}
		}
	}

	go connection.DeleteRedis(RedisAdminAllKey)
	return nil
}
