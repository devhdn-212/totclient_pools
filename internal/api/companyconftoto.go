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

type companyconftotoApi struct {
	companyconftotoService domain.CompanyconftotoService
}

func NewCompanyconftotoApi(app *fiber.App,
	companyconftotoService domain.CompanyconftotoService,
	authmidle fiber.Handler) {
	ad := companyconftotoApi{
		companyconftotoService: companyconftotoService,
	}
	companyconftoto := app.Group("/api/companyconftoto", authmidle)
	companyconftoto.Post("", ad.Index)
	companyconftoto.Post("/save", ad.Save)
}
func (co *companyconftotoApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CompanyconftotoAll
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
		connection.Log.Warn("Validation failed for Company Conf TOTO All",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	flagpage := util.Validpage(client_username, "COMPANYCONFTOTO-VIEW")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}
	res, err := co.companyconftotoService.All(c, req.IDcompany)
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
func (co *companyconftotoApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CompanyconftotoSave
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
		connection.Log.Warn("Validation failed for update Company Conf TOTO",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}
	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	flagpage := util.Validpage(client_username, "COMPANYCONFTOTO-SAVE")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}

	err := co.companyconftotoService.Save(c, req, client_username)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update Company Conf TOTO",
			zap.String("id", req.IDcompconftoto),
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
	connection.Log.Info("Company Conf TOTO create / update successfully",
		zap.String("id", req.IDcompconftoto),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
