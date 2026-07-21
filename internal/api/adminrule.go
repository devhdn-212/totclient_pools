package api

import (
	"context"
	"net/http"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"

	"github.com/gofiber/fiber/v2"
)

type adminruleApi struct {
	adminruleService domain.AdminruleService
}

func NewAdminruleApi(app *fiber.App,
	adminruleService domain.AdminruleService,
	authmidle fiber.Handler) {
	ad := adminruleApi{
		adminruleService: adminruleService,
	}
	admin := app.Group("/api/adminrule", authmidle)
	admin.Post("", ad.Index)
}
func (ad *adminruleApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ad.adminruleService.All(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.JSON(dto.CreateResponseSuccess(res))
}
