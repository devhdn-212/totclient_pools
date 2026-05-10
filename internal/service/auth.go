package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/dto"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"
	"github.com/devhdn-212/gofibergoqu_master/internal/connection"
	"github.com/devhdn-212/gofibergoqu_master/internal/repository"
	"github.com/devhdn-212/gofibergoqu_master/internal/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	RedisClient = "master:client:"
)

type authService struct {
	db                  *pgxpool.Pool
	conf                *config.Config
	adminRepository     domain.AdminsRepository
	adminruleRepository domain.AdminruleRepository
}

func NewAuth(db *pgxpool.Pool,
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
	// 1. Cari user berdasarkan username
	user, err := a.adminRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user.Username == "" {
		return dto.AuthResponse{}, errors.New("Username / Password Not Found")
	}

	// 2. Verifikasi Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("Username / Password Not Found")
	}

	// 3. Ambil Rule (Hak Akses)
	rule, errrule := a.adminruleRepository.GetRule(ctx, user.Idadmin)
	if errrule != nil {
		return dto.AuthResponse{}, errors.New("Please contact Admin")
	}

	// 4. Update Last Login menggunakan Transaksi PGX
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	defer tx.Rollback(ctx)

	txExec := repository.NewPGXTxExecutor(tx)
	txRepo := repository.NewAdminRepository(txExec)

	// Update data login
	user.Ipaddress = req.Ipaddress
	user.Lastlogin = sql.NullTime{Valid: true, Time: util.GetNowJakarta()}

	// Gunakan txRepo agar masuk dalam scope transaksi
	if err = txRepo.UpdateLogin(ctx, &user); err != nil {
		return dto.AuthResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.AuthResponse{}, err
	}

	// 5. Simpan Data ke Redis
	var clientRedis dto.AuthClientRedis
	clientRedis.Username = user.Username
	clientRedis.IDrule = user.Idadmin
	clientRedis.Rule = rule

	// Enkripsi username untuk payload token
	dataclient_encr, keymap := util.Encryption(user.Username)
	dataclient_encr_final := dataclient_encr + "|" + strconv.Itoa(keymap)

	go connection.SetRedis(RedisClient+user.Username, clientRedis, 1440*time.Minute)

	// 6. Generate JWT Token (v5)
	claim := jwt.MapClaims{
		"username":    user.Username, // Tambahkan username plain untuk middleware
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
