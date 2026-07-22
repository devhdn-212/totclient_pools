package domain

import (
	"context"
	"errors"

	"github.com/devhdn-212/totclient_api/dto"
)

// ErrInsufficientBalance is returned when the player's balance (from the
// pusat/wallet balance check) is lower than the checkout's declared total —
// checked before anything is written, so a short balance never touches the
// DB at all.
var ErrInsufficientBalance = errors.New("insufficient balance")

// ErrPasaranOffline is returned when the target pasaran's schedule marks it
// as OFFLINE (closed for the day) at checkout time — same "reject before
// anything is written" treatment as ErrInsufficientBalance.
var ErrPasaranOffline = errors.New("pasaran offline")

type CheckoutService interface {
	Submit(ctx context.Context, req dto.CheckoutRequest, ipaddress string) (dto.CheckoutResponse, error)
}
