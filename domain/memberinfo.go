package domain

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/dto"
)

// ErrInvalidToken is returned when the token is not recognized by this
// (pusat) server, so callers can distinguish it from other error sources.
var ErrInvalidToken = errors.New("invalid token")

type Memberinfo struct {
	Username  string
	Balance   int64
	AgentCode string
	Token     string
}

type MemberinfoService interface {
	CheckToken(ctx context.Context, rec dto.MemberinfoResponse) (*Memberinfo, error)
}
