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

type resultApi struct {
	resultService domain.ResultService
}

// NewResultApi — public route, no JWT authmidle, same as /api/transaksi:
// player identity is the raw launch token, re-checked against the
// pusat/wallet server on every call. Draw results themselves aren't
// player-specific, but the token check still gates the endpoint to
// sessions with a live launch token.
func NewResultApi(app *fiber.App, resultService domain.ResultService) {
	aa := &resultApi{
		resultService: resultService,
	}
	result := app.Group("/api/result")
	result.Post("", aa.Fetch)
}

func (a resultApi) Fetch(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.ResultRequest
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Result"),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for Result",
			zap.Any("validation_errors", fails),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	res, err := a.resultService.Fetch(c, req)
	if errors.Is(err, domain.ErrInvalidToken) {
		return ctx.Status(http.StatusUnauthorized).
			JSON(dto.CreateResponseErrorCode(http.StatusUnauthorized, dto.ErrCodeInvalidTokenPusat, "invalid token"))
	}
	if errors.Is(err, util.ErrNotFound) {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseError(http.StatusBadRequest, "pasaran tidak ditemukan"))
	}
	if err != nil {
		connection.Log.Error("Result failed", zap.String("error", err.Error()))
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseErrorCode(http.StatusInternalServerError, dto.ErrCodeInternalLocal, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(res))
}
