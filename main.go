package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devhdn-212/totagen_api/dto"
	"github.com/devhdn-212/totagen_api/internal/api"
	"github.com/devhdn-212/totagen_api/internal/config"
	"github.com/devhdn-212/totagen_api/internal/connection"
	"github.com/devhdn-212/totagen_api/internal/repository"
	"github.com/devhdn-212/totagen_api/internal/service"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := NewGCPLogger()

	connection.SetLogger(logger)

	cnf := config.Get()
	// 3. Koneksi Database (pgxpool)
	dbPool := connection.GetDatabase(cnf.Database)
	defer dbPool.Close()

	// 4. Inisialisasi & Health Check Redis
	if err := connection.InitRedis(cnf.Redis); err != nil {
		logger.Fatal("Failed to init Redis", zap.Error(err))
	}
	defer connection.RDB.Close()

	if !connection.RedisHealth() {
		logger.Fatal("Redis is not healthy")
	}

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(etag.New())
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		rid, _ := c.Locals("requestid").(string)
		fields := []zap.Field{
			zap.String("request_id", rid),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int64("latency_ms", latency.Microseconds()),
			zap.String("ip", c.IP()),
		}
		// Optional: POST JSON (HATI-HATI DI PROD)
		if c.Method() == fiber.MethodPost && c.Is("json") {
			if len(c.Body()) < 2048 { // guard
				var parsed map[string]interface{}
				if err := json.Unmarshal(c.Body(), &parsed); err == nil {
					delete(parsed, "password")
					delete(parsed, "token")
					fields = append(fields, zap.Any("json_body", parsed))
				}
			}
		}

		// log
		logger.Info("http_request", fields...)
		return err
	})

	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key), JWTAlg: "HS256"},
		SuccessHandler: func(c *fiber.Ctx) error {
			token, ok := c.Locals("user").(*jwt.Token)
			if !ok || token == nil {
				return c.Status(fiber.StatusUnauthorized).
					JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).
					JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
			}
			username, ok := claims["clien_admin"].(string)
			c.Locals("client_username", username)
			if !ok || username == "" {
				return c.Status(fiber.StatusUnauthorized).
					JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token - Username"))
			}
			jti, ok := claims["jti"].(string)
			if !ok || jti == "" {
				return c.Status(fiber.StatusUnauthorized).
					JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
			}
			isBlacklisted, err := connection.IsJWTBlacklisted(jti)
			if err != nil || isBlacklisted {
				return c.Status(fiber.StatusUnauthorized).
					JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
			}
			if !validateJwtClaims(claims, cnf.Jwt.Issuer, cnf.Jwt.Audience) {
				return c.Status(fiber.StatusUnauthorized).
					JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "invalid token"))
			}
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).
				JSON(dto.CreateResponseError(fiber.StatusUnauthorized, "missing token, please login"))
		},
	})
	pgxExec := repository.NewPGXExecutor(dbPool)
	adminRepository := repository.NewAdminRepository(pgxExec)
	adminruleRepository := repository.NewAdminruleRepository(pgxExec)
	clientruleRepository := repository.NewClientruleRepository(pgxExec)
	currRepository := repository.NewCurrRepository(pgxExec)
	domainRepository := repository.NewDomainRepository(pgxExec)
	groupcompanyRepository := repository.NewGroupcompanyRepository(pgxExec)
	companyRepository := repository.NewCompanyRepository(pgxExec)
	companywalletRepository := repository.NewCompanywalletRepository(pgxExec)
	companyadminRepository := repository.NewCompanyadminRepository(pgxExec)
	companypasaranRepository := repository.NewCompanypasaranRepository(pgxExec)
	companyconftotoRepository := repository.NewCompanyconftotoRepository(pgxExec)
	customerRepository := repository.NewCustomerRepository(pgxExec)
	pasarantotoRepository := repository.NewPasarantotoRepository(pgxExec)

	adminService := service.NewAdminService(dbPool, adminRepository)
	adminruleService := service.NewAdminruleService(dbPool, adminruleRepository)
	clinetruleService := service.NewClientruleService(dbPool, clientruleRepository)
	currService := service.NewCurrService(dbPool, currRepository)
	domainService := service.NewDomainService(dbPool, domainRepository)
	groupcompanyService := service.NewGroupCompanyService(dbPool, groupcompanyRepository)
	companyService := service.NewCompanyService(dbPool, companyRepository)
	companywalletService := service.NewCompanywalletService(dbPool, companywalletRepository)
	companyadminService := service.NewCompanyadminService(dbPool, companyadminRepository)
	companypasaranService := service.NewCompaypasaranService(dbPool, companypasaranRepository)
	companyconftotoService := service.NewCompanyconftotoService(dbPool, companyconftotoRepository)
	customerService := service.NewCustomerService(dbPool, customerRepository)
	pasarantotoService := service.NewPasarantotoService(dbPool, pasarantotoRepository)
	authService := service.NewAuth(dbPool, cnf, adminRepository, adminruleRepository)

	api.NewAdminApi(app, adminService, adminruleService, jwtMidd)
	api.NewAdminruleApi(app, adminruleService, jwtMidd)
	api.NewClientruleApi(app, clinetruleService, jwtMidd)
	api.NewCurrApi(app, currService, jwtMidd)
	api.NewDomainApi(app, domainService, jwtMidd)
	api.NewGroupcompanyApi(app, groupcompanyService, jwtMidd)
	api.NewCompanyApi(app, companyService, groupcompanyService, currService, clinetruleService, pasarantotoService, jwtMidd)
	api.NewCompanywalletApi(app, companywalletService, jwtMidd)
	api.NewCompanyadminApi(app, companyadminService, jwtMidd)
	api.NewCompanypasaranApi(app, companypasaranService, jwtMidd)
	api.NewCompanyconftotoApi(app, companyconftotoService, jwtMidd)
	api.NewPasarantotoApi(app, pasarantotoService, jwtMidd)
	api.NewCustomer(app, customerService, jwtMidd)
	api.NewAuth(app, authService, jwtMidd)

	go func() {
		appsPort := cnf.Server.Port
		if err := app.Listen(":" + appsPort); err != nil {
			logger.Fatal("Failed to start app", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Gracefully shutting down...")
	_ = app.Shutdown()

	logger.Info("Running cleanup tasks...")

	// Your cleanup tasks go here
	dbPool.Close()
	// Berikan timeout untuk shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Shutdown complete")
}

func validateJwtClaims(claims jwt.MapClaims, issuer, audience string) bool {
	if issuer != "" {
		iss, ok := claims["iss"].(string)
		if !ok || iss != issuer {
			return false
		}
	}
	if audience != "" {
		switch aud := claims["aud"].(type) {
		case string:
			if aud != audience {
				return false
			}
		case []interface{}:
			found := false
			for _, v := range aud {
				if s, ok := v.(string); ok && s == audience {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case []string:
			found := false
			for _, s := range aud {
				if s == audience {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func NewGCPLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05"))
	}
	encoderConfig.LevelKey = "severity"
	encoderConfig.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	return zap.New(core, zap.AddCaller())
}
