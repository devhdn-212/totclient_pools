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

type companyApi struct {
	companyService      domain.CompanyService
	groupcompanyService domain.GroupcompanyService
	currService         domain.CurrencyService
	clientruleService   domain.ClientruleService
	pasarantotoService  domain.PasarantotoService
}

func NewCompanyApi(app *fiber.App,
	companyService domain.CompanyService,
	groupcompanyService domain.GroupcompanyService,
	currService domain.CurrencyService,
	clientruleService domain.ClientruleService,
	pasarantotoService domain.PasarantotoService,
	authmidle fiber.Handler) {
	ad := companyApi{
		companyService:      companyService,
		groupcompanyService: groupcompanyService,
		currService:         currService,
		clientruleService:   clientruleService,
		pasarantotoService:  pasarantotoService,
	}
	company := app.Group("/api/company", authmidle)
	company.Post("", ad.Index)
	company.Post("/save", ad.Save)
}
func (co *companyApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	resselectgroup, errselectgroup := co.groupcompanyService.Select(c)
	if errselectgroup != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error - Group Company"))
	}

	resselectpasaran, errselectpasaran := co.pasarantotoService.Select(c)
	if errselectpasaran != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error - Pasarantoto"))
	}

	resselectrule, errselectrule := co.clientruleService.Select(c)
	if errselectrule != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error - Rule"))
	}

	resselect, errselect := co.currService.Select(c)
	if errselect != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error - Currency"))
	}
	res, err := co.companyService.All(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.JSON(fiber.Map{
		"status":          fiber.StatusOK,
		"message":         "success",
		"listpasarantoto": resselectpasaran,
		"listcurr":        resselect,
		"listrule":        resselectrule,
		"listgroup":       resselectgroup,
		"record":          res,
	})
}
func (co *companyApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CompanySave
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
		connection.Log.Warn("Validation failed for update Company",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}
	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	flagpage := util.Validpage(client_username, "COMPANY-SAVE")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}

	err := co.companyService.Save(c, req, client_username)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update Company",
			zap.String("id", req.ID),
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
	connection.Log.Info("Company create / update successfully",
		zap.String("id", req.ID),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
