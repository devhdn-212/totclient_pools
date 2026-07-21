package api

import (
	"context"
	"net/http"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/util"

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
