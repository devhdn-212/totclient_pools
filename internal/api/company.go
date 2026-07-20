package api

import (
	"context"
	"net/http"
	"time"

	"github.com/devhdn-212/totagen_api/domain"
	"github.com/devhdn-212/totagen_api/dto"

	"github.com/gofiber/fiber/v2"
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
