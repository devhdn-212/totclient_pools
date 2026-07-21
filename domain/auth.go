package domain

import (
	"context"

	"github.com/devhdn-212/totclient_api/dto"
)

type AuthService interface {
	Login(ctx context.Context, red dto.AuthRequest) (dto.AuthResponse, error)
}
