package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/util"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type checkoutApi struct {
	checkoutService domain.CheckoutService
	settingService  domain.SettingService
}

// NewCheckoutApi — public route, no JWT authmidle. Player identity/session
// is the raw launch token itself (checked against the pusat/wallet server
// inside the service), same as /api/servicetoken — there's no separate
// player JWT in this system.
func NewCheckoutApi(app *fiber.App, checkoutService domain.CheckoutService, settingService domain.SettingService) {
	aa := &checkoutApi{
		checkoutService: checkoutService,
		settingService:  settingService,
	}
	checkout := app.Group("/api/checkout")
	checkout.Post("", aa.Submit)
}

func (a checkoutApi) Submit(ctx *fiber.Ctx) error {
	// Chunked baskets can run to hundreds of rows per chunk, each costing up
	// to ~4 sequential Redis round-trips (limittotal + limitglobal, each an
	// EXISTS check plus a Lua script call against a remote Redis) — give
	// this more room than the usual 10s endpoint timeout.
	c, cancel := context.WithTimeout(ctx.Context(), 60*time.Second)
	defer cancel()

	// Server-side gate — the frontend already blocks play off the
	// /api/servicetoken alert, but a tab left open from before maintenance
	// started must not be able to slip a bet through anyway.
	if handled, err := respondIfMaintenance(ctx, c, a.settingService, "checkout"); handled {
		return err
	}

	var req dto.CheckoutRequest
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Checkout"),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for Checkout",
			zap.Any("validation_errors", fails),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	ipaddress := req.Ipaddress
	if ipaddress == "" {
		ipaddress = ctx.IP()
	}

	res, err := a.checkoutService.Submit(c, req, ipaddress)
	if errors.Is(err, domain.ErrInvalidToken) {
		return ctx.Status(http.StatusUnauthorized).
			JSON(dto.CreateResponseErrorCode(http.StatusUnauthorized, dto.ErrCodeInvalidTokenPusat, "invalid token"))
	}
	if errors.Is(err, domain.ErrInsufficientBalance) {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorCode(http.StatusBadRequest, dto.ErrCodeInsufficientBalancePusat, "insufficient balance"))
	}
	if errors.Is(err, domain.ErrPasaranOffline) {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorCode(http.StatusBadRequest, dto.ErrCodePasaranOfflinePusat, "pasaran sedang tutup"))
	}
	if errors.Is(err, util.ErrNotFound) {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseError(http.StatusBadRequest, "pasaran tidak ditemukan"))
	}
	if err != nil {
		connection.Log.Error("Checkout failed: " + err.Error())
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseErrorCode(http.StatusInternalServerError, dto.ErrCodeInternalLocal, dto.MsgInternalError))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(res))
}
