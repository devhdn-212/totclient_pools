package domain

import (
	"context"
	"gofibergocu/dto"
)

type AuthService interface {
	Login(ctx context.Context, red dto.AuthRequest) (dto.AuthResponse, error)
}
