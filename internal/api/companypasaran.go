package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/devhdn-212/totmaster_api/domain"
	"github.com/devhdn-212/totmaster_api/dto"
	"github.com/devhdn-212/totmaster_api/internal/connection"
	"github.com/devhdn-212/totmaster_api/internal/util"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type companypasaranApi struct {
	companypasaaranService domain.CompanypasaranService
}

func NewCompanypasaranApi(app *fiber.App,
	companypasaranService domain.CompanypasaranService,
	authmidle fiber.Handler) {
	ad := companypasaranApi{
		companypasaaranService: companypasaranService,
	}
	companypasaran := app.Group("/api/companypasaran", authmidle)
	companypasaran.Post("", ad.Index)
	companypasaran.Post("/save", ad.Save)
}
func (co *companypasaranApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CompanypasaranAll
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Company Admin All"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for Company Pasaran All",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	flagpage := util.Validpage(client_username, "COMPANYADMIN-VIEW")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}
	res, err := co.companypasaaranService.All(c, req.IDcompany)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"record":  res,
	})
}
func (co *companypasaranApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CompanypasaranSave
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Create Company Admin"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for update Company Pasaran",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}
	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	flagpage := util.Validpage(client_username, "COMPANYPASARAN-SAVE")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}

	err := co.companypasaaranService.Save(c, req, client_username)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update Company Pasaran",
			zap.String("id", req.IDcomppasaran),
			zap.String("error", err.Error()),
			zap.String("record", string(recordJson)),
		)
		// cek duplicate entry
		if err.Error() == "duplicate entry" {
			return ctx.Status(http.StatusConflict).
				JSON(dto.CreateResponseError(http.StatusConflict, "Duplicate Entry"))
		}
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	connection.Log.Info("Company Pasaran create / update successfully",
		zap.String("id", req.IDcomppasaran),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
