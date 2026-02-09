package api

import (
	"context"
	"encoding/json"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type adminApi struct {
	adminService     domain.AdminService
	adminruleService domain.AdminruleService
}

func NewAdminApi(app *fiber.App,
	adminService domain.AdminService, adminruleService domain.AdminruleService,
	authmidle fiber.Handler) {
	ad := adminApi{
		adminService:     adminService,
		adminruleService: adminruleService,
	}
	admin := app.Group("/api/admin", authmidle)
	admin.Post("", ad.Index)
	admin.Post("/save", ad.Save)
}
func (ad *adminApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	resselect, errselect := ad.adminruleService.Select(c)
	if errselect != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	res, err := ad.adminService.All(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.JSON(fiber.Map{
		"status":        fiber.StatusOK,
		"message":       "success",
		"listadminrule": resselect,
		"record":        res,
	})
}
func (ad *adminApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AdminSave
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Create Admin"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
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
	datatoken := ctx.Locals("client_username").(string)
	client_username, _ := util.Parsing_final(datatoken)

	err := ad.adminService.Save(c, req, client_username)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update admin",
			zap.String("id", req.Username),
			zap.String("error", err.Error()),
			zap.String("record", string(recordJson)),
		)
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	connection.Log.Info("Admin create / update successfully",
		zap.String("id", req.Username),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
