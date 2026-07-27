package domain

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/dto"
)

// ErrInvalidToken is returned when the token is not recognized by this
// (pusat) server, so callers can distinguish it from other error sources.
var ErrInvalidToken = errors.New("invalid token")

// ErrInvalidAgent / ErrInvalidMarket are returned when the launch params'
// agent or market code has no matching (active) row in the database — a
// permanent rejection, not a transient failure. Checked in that order:
// market is looked up scoped to the agent, so an unknown agent must be
// rejected first or it would otherwise be misreported as an unknown market.
var ErrInvalidAgent = errors.New("invalid agent")
var ErrInvalidMarket = errors.New("invalid market")

type Memberinfo struct {
	Username  string
	Balance   int64
	AgentCode string
	Token     string
}

type MemberinfoService interface {
	CheckToken(ctx context.Context, rec dto.MemberinfoResponse) (*Memberinfo, error)
}
