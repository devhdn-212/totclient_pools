package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/util"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type pasarantotoApi struct {
	pasarantotoService domain.PasarantotoService
}

func NewPasarantotoApi(app *fiber.App,
	pasarantotoService domain.PasarantotoService,
	authmidle fiber.Handler) {
	ad := pasarantotoApi{
		pasarantotoService: pasarantotoService,
	}
	pasarantoto := app.Group("/api/pasarantoto", authmidle)
	pasarantoto.Post("", ad.Index)
	pasarantoto.Post("/save", ad.Save)
}
func (ad *pasarantotoApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ad.pasarantotoService.All(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.JSON(dto.CreateResponseSuccess(res))
}
func (ad *pasarantotoApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.PasarantotoSave
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Create Pasarantoto"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for update pasarantoto",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}
	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	flagpage := util.Validpage(client_username, "PASARANTOTO-SAVE")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}

	err := ad.pasarantotoService.Save(c, req, client_username)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update Pasarantoto",
			zap.String("id", req.IDpasarantogel),
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
	connection.Log.Info("Pasarantoto create / update successfully",
		zap.String("id", req.IDpasarantogel),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
