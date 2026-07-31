package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/dto"
	"github.com/devhdn-212/totclient_pools/internal/connection"
	"github.com/devhdn-212/totclient_pools/internal/util"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type memberinfoApi struct {
	memberinfoService domain.MemberinfoService
	settingService    domain.SettingService
}

func NewMemberInfo(app *fiber.App,
	memberinfoService domain.MemberinfoService,
	settingService domain.SettingService) {
	aa := &memberinfoApi{
		memberinfoService: memberinfoService,
		settingService:    settingService,
	}
	memberinfo := app.Group("/api/servicetoken")
	memberinfo.Post("", aa.Checktoken)
}
func (a memberinfoApi) Checktoken(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	if handled, err := respondIfMaintenance(ctx, c, a.settingService, "servicetoken"); handled {
		return err
	}

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

	res, err := a.memberinfoService.CheckToken(c, req)
	if errors.Is(err, domain.ErrInvalidToken) {
		return ctx.Status(http.StatusUnauthorized).
			JSON(dto.CreateResponseErrorCode(http.StatusUnauthorized, dto.ErrCodeInvalidTokenPusat, "invalid token"))
	}
	if errors.Is(err, domain.ErrInvalidAgent) {
		return ctx.Status(http.StatusNotFound).
			JSON(dto.CreateResponseErrorCode(http.StatusNotFound, dto.ErrCodeInvalidAgentPusat, "Agent tidak ditemukan. Please Contact Admin."))
	}
	if errors.Is(err, domain.ErrInvalidMarket) {
		return ctx.Status(http.StatusNotFound).
			JSON(dto.CreateResponseErrorCode(http.StatusNotFound, dto.ErrCodeInvalidMarketPusat, "Market tidak ditemukan. Please Contact Admin."))
	}
	if err != nil {
		connection.Log.Error("CheckToken failed: " + err.Error())
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseErrorCode(http.StatusInternalServerError, dto.ErrCodeInternalLocal, dto.MsgInternalError))
	}
	return ctx.JSON(fiber.Map{
		"status":   fiber.StatusOK,
		"message":  "success",
		"username": res.Username,
		"balance":  res.Balance,
		"token":    res.Token,
	})
}
