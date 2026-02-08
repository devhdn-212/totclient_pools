package service

import (
	"context"
	"errors"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	conf            *config.Config
	adminRepository domain.AdminsRepository
}

func NewAuth(cnf *config.Config,
	adminRepository domain.AdminsRepository) domain.AuthService {
	return authService{
		conf:            cnf,
		adminRepository: adminRepository,
	}
}
func (a authService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.adminRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user.Username == "" {
		return dto.AuthResponse{}, errors.New("auth failed")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("auth failed")
	}

	claim := jwt.MapClaims{
		"username": user.Username,
		"jti":      uuid.NewString(),
		"iss":      a.conf.Jwt.Issuer,
		"aud":      a.conf.Jwt.Audience,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenstr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("auth failed")
	}
	return dto.AuthResponse{Token: tokenstr}, nil

}
