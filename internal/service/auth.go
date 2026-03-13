package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/dto"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"
	"github.com/devhdn-212/gofibergoqu_master/internal/connection"
	"github.com/devhdn-212/gofibergoqu_master/internal/repository"
	"github.com/devhdn-212/gofibergoqu_master/internal/util"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	RedisClient = "master:client:"
)

type authService struct {
	db                  *sql.DB
	conf                *config.Config
	adminRepository     domain.AdminsRepository
	adminruleRepository domain.AdminruleRepository
}

func NewAuth(db *sql.DB,
	cnf *config.Config,
	adminRepository domain.AdminsRepository,
	adminruleRepository domain.AdminruleRepository) domain.AuthService {
	return authService{
		db:                  db,
		conf:                cnf,
		adminRepository:     adminRepository,
		adminruleRepository: adminruleRepository,
	}
}
func (a authService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.adminRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user.Username == "" {
		return dto.AuthResponse{}, errors.New("Username / Password Not Found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("Username / Password Not Found")
	}

	rule, errrule := a.adminruleRepository.GetRule(ctx, user.Idadmin)
	if errrule != nil {
		return dto.AuthResponse{}, errors.New("Please contact Admin")
	}

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	txExec := repository.NewGoquTxExecutor(tx)
	txRepo := repository.NewAdminRepository(txExec)
	flagupdate, errupdate := txRepo.FindByUsername(ctx, user.Username)
	if errupdate != nil {
		return dto.AuthResponse{}, errupdate
	}
	flagupdate.Username = user.Username
	flagupdate.Ipaddress = req.Ipaddress
	flagupdate.Lastlogin = sql.NullTime{Valid: true, Time: time.Now().In(loc)}
	if err = a.adminRepository.UpdateLogin(ctx, &flagupdate); err != nil {
		fmt.Println(err)
		return dto.AuthResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return dto.AuthResponse{}, err
	}

	var clientRedis dto.AuthClientRedis
	clientRedis.Username = user.Username
	clientRedis.IDrule = user.Idadmin
	clientRedis.Rule = rule

	dataclient := user.Username
	dataclient_encr, keymap := util.Encryption(dataclient)
	dataclient_encr_final := dataclient_encr + "|" + strconv.Itoa(keymap)

	go connection.SetRedis(RedisClient+user.Username, clientRedis, 1440*time.Minute)

	claim := jwt.MapClaims{
		"clien_admin": dataclient_encr_final,
		"jti":         uuid.NewString(),
		"iss":         a.conf.Jwt.Issuer,
		"aud":         a.conf.Jwt.Audience,
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenstr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("auth failed")
	}
	return dto.AuthResponse{Token: tokenstr}, nil

}
