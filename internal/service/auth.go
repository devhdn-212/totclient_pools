package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/config"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authService struct {
	conf           *config.Config
	userRepository domain.UserRepository
}

func NewAuth(cnf *config.Config,
	userRepository domain.UserRepository) domain.AuthService {
	return authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}
func (a authService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user.ID == "" {
		return dto.AuthResponse{}, errors.New("auth failed")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("auth failed")
	}

	claim := jwt.MapClaims{
		"id":  user.ID,
		"jti": uuid.NewString(),
		"iss": a.conf.Jwt.Issuer,
		"aud": a.conf.Jwt.Audience,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenstr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("auth failed")
	}
	return dto.AuthResponse{Token: tokenstr}, nil

}
