package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/domain"
	"github.com/devhdn-212/gofibergoqu_master/dto"
	"github.com/devhdn-212/gofibergoqu_master/internal/connection"
	"github.com/devhdn-212/gofibergoqu_master/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type authApi struct {
	authService      domain.AuthService
	adminruleService domain.AdminruleService
}

func NewAuth(app *fiber.App,
	authService domain.AuthService,
	authmidle fiber.Handler) {
	aa := &authApi{
		authService: authService,
	}
	auth := app.Group("/api/auth", limiter.New(limiter.Config{
		Max:        20,
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
	auth.Post("/page", authmidle, aa.Page)
	auth.Post("/logout", authmidle, aa.Logout)
}
func (a authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Login"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for update Admin",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	res, err := a.authService.Login(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"token":   res.Token,
	})
}

func (a authApi) Page(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthPage
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Login"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for update Admin",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	datatoken := ctx.Locals("client_username").(string)
	client_username := util.Parsing_final(datatoken)
	fmt.Println("Username: ", client_username)
	flagpage := util.Validpage(client_username, req.Page)
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}
	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
	})
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
	username, ok := claims["clien_admin"].(string)
	ctx.Locals("client_username", username)
	client_username := util.Parsing_final(username)
	fmt.Println("USERNAME: ", client_username)
	if !ok || username == "" {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
	} else {
		go connection.DeleteRedis("master:client:" + client_username)
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
	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
	})
}
