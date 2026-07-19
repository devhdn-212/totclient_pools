package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/util"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type pasaranApi struct {
	pasaranService domain.PasaranService
}

func NewPasaranApi(app *fiber.App,
	pasaranService domain.PasaranService,
	authmidle fiber.Handler) {
	ad := pasaranApi{
		pasaranService: pasaranService,
	}
	pasaran := app.Group("/api/pasaran", authmidle)
	pasaran.Post("", ad.Index)
	pasaran.Post("/save", ad.Save)
}
func (co *pasaranApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	datatoken := ctx.Locals("client_agen").(string)
	client_username := util.Parsing_final(datatoken)
	flag, _, idcomp := util.GetDataRedisClient(client_username)
	fmt.Println("COMP : ", idcomp)
	if !flag {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	flagpage := util.Validpage(client_username, "COMPANY-VIEW")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}
	res, err := co.pasaranService.All(c, idcomp)
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
func (co *pasaranApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.PasaranSave
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
	datatoken := ctx.Locals("client_agen").(string)
	client_username := util.Parsing_final(datatoken)
	flag, _, idcomp := util.GetDataRedisClient(client_username)
	fmt.Println("COMP : ", idcomp)
	if !flag {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	flagpage := util.Validpage(client_username, "COMPANY-SAVE")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}

	err := co.pasaranService.Save(c, req, client_username, idcomp)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update Pasaran",
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
	connection.Log.Info("Pasaran create / update successfully",
		zap.String("id", req.IDcomppasaran),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
