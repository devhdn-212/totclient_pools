package api

import (
	"context"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/connection"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v5"
)

type authApi struct {
	authService domain.AuthService
}

func NewAuth(app *fiber.App, authService domain.AuthService, authmidle fiber.Handler) {
	aa := &authApi{
		authService: authService,
	}
	auth := app.Group("/auth", limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(dto.CreateResponseError(fiber.StatusTooManyRequests, "too many requests"))
		},
	}))
	auth.Post("", aa.Login)
	auth.Post("/logout", authmidle, aa.Logout)
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
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	ctx.Locals("client_admin", req.Username)
	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(res))
}

func (a authApi) Logout(ctx *fiber.Ctx) error {
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
	}
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError(fiber.StatusBadRequest, "missing token id"))
	}
	expUnix, ok := claims["exp"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError(fiber.StatusBadRequest, "missing token expiry"))
	}
	ttl := time.Until(time.Unix(int64(expUnix), 0))
	if ttl <= 0 {
		return ctx.SendStatus(fiber.StatusNoContent)
	}
	if err := connection.BlacklistJWT(jti, ttl); err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
