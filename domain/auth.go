package domain

import (
	"context"

	"github.com/devhdn-212/totmaster_api/dto"
)

type AuthService interface {
	Login(ctx context.Context, red dto.AuthRequest) (dto.AuthResponse, error)
}
