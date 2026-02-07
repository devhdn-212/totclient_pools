package api

import (
	"context"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authApi struct {
	authService domain.AuthService
}

func NewAuth(app *fiber.App, authService domain.AuthService) {
	aa := &authApi{
		authService: authService,
	}
	app.Post("/auth", aa.Login)
}
func (a authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	res, err := a.authService.Login(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusBadRequest, err.Error()))
	}
	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(res))
}
