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
}

// NewCheckoutApi — public route, no JWT authmidle. Player identity/session
// is the raw launch token itself (checked against the pusat/wallet server
// inside the service), same as /api/servicetoken — there's no separate
// player JWT in this system.
func NewCheckoutApi(app *fiber.App, checkoutService domain.CheckoutService) {
	aa := &checkoutApi{
		checkoutService: checkoutService,
	}
	checkout := app.Group("/api/checkout")
	checkout.Post("", aa.Submit)
}

func (a checkoutApi) Submit(ctx *fiber.Ctx) error {
	// Chunked baskets can run to thousands of rows — give this more room
	// than the usual 10s endpoint timeout.
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

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
	if errors.Is(err, util.ErrNotFound) {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseError(http.StatusBadRequest, "pasaran tidak ditemukan"))
	}
	if err != nil {
		connection.Log.Error("Checkout failed", zap.String("error", err.Error()))
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseErrorCode(http.StatusInternalServerError, dto.ErrCodeInternalLocal, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(res))
}
