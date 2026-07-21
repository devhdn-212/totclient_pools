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
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

type serviceinitApi struct {
	pasaranService domain.PasaranService
}

func NewServiceinit(app *fiber.App,
	pasaranService domain.PasaranService) {
	aa := &serviceinitApi{
		pasaranService: pasaranService,
	}
	serviceinit := app.Group("/api/serviceinit", limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(dto.CreateResponseError(fiber.StatusTooManyRequests, "too many requests"))
		},
	}))
	serviceinit.Post("", aa.Init)
}
func (a serviceinitApi) Init(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.MemberinfoResponse
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Token"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for update Admin",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	res, err := a.pasaranService.FindID(c, req.Agen, req.Market)
	if errors.Is(err, domain.ErrInvalidToken) {
		return ctx.Status(http.StatusUnauthorized).
			JSON(dto.CreateResponseErrorCode(http.StatusUnauthorized, dto.ErrCodeInvalidTokenPusat, "invalid token"))
	}
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseErrorCode(http.StatusInternalServerError, dto.ErrCodeInternalLocal, err.Error()))
	}
	return ctx.JSON(fiber.Map{
		"status":      fiber.StatusOK,
		"message":     "success",
		"initpasaran": res,
	})
}
