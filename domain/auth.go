package domain

import (
	"context"

	"github.com/devhdn-212/gofibergoqu_master/dto"
)

type AuthService interface {
	Login(ctx context.Context, red dto.AuthRequest) (dto.AuthResponse, error)
}
