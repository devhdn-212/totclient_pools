package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devhdn-212/gofibergoqu_master/dto"
	"github.com/devhdn-212/gofibergoqu_master/internal/api"
	"github.com/devhdn-212/gofibergoqu_master/internal/config"
	"github.com/devhdn-212/gofibergoqu_master/internal/connection"
	"github.com/devhdn-212/gofibergoqu_master/internal/repository"
	"github.com/devhdn-212/gofibergoqu_master/internal/service"

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
	dbConnection := connection.GetDatabase(cnf.Database)
	if err := connection.InitRedis(cnf.Redis); err != nil {
		panic(err)
	}
	defer connection.RDB.Close()
	if !connection.RedisHealth() {
		panic("Redis is not healthy")
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
	goquExec := repository.NewGoquExecutor(dbConnection)
	adminRepository := repository.NewAdminRepository(goquExec)
	adminruleRepository := repository.NewAdminruleRepository(goquExec)
	clientruleRepository := repository.NewClientruleRepository(goquExec)
	currRepository := repository.NewCurrRepository(goquExec)
	uomRepository := repository.NewUomRepository(goquExec)
	bankRepository := repository.NewBankRepository(goquExec)
	domainRepository := repository.NewDomainRepository(goquExec)
	companyRepository := repository.NewCompanyRepository(goquExec)
	companywalletRepository := repository.NewCompanywalletRepository(goquExec)
	companyadminRepository := repository.NewCompanyadminRepository(goquExec)
	customerRepository := repository.NewCustomerRepository(goquExec)

	adminService := service.NewAdminService(dbConnection, adminRepository)
	adminruleService := service.NewAdminruleService(dbConnection, adminruleRepository)
	clinetruleService := service.NewClientruleService(dbConnection, clientruleRepository)
	currService := service.NewCurrService(dbConnection, currRepository)
	uomService := service.NewUomService(dbConnection, uomRepository)
	bankService := service.NewBankService(dbConnection, bankRepository)
	domainService := service.NewDomainService(dbConnection, domainRepository)
	companyService := service.NewCompanyService(dbConnection, companyRepository)
	companywalletService := service.NewCompanywalletService(dbConnection, companywalletRepository)
	companyadminService := service.NewCompanyadminService(dbConnection, companyadminRepository)
	customerService := service.NewCustomerService(dbConnection, customerRepository)
	authService := service.NewAuth(dbConnection, cnf, adminRepository, adminruleRepository)

	api.NewAdminApi(app, adminService, adminruleService, jwtMidd)
	api.NewAdminruleApi(app, adminruleService, jwtMidd)
	api.NewClientruleApi(app, clinetruleService, jwtMidd)
	api.NewCurrApi(app, currService, jwtMidd)
	api.NewUomApi(app, uomService, jwtMidd)
	api.NewBankApi(app, bankService, jwtMidd)
	api.NewDomainApi(app, domainService, jwtMidd)
	api.NewCompanyApi(app, companyService, currService, clinetruleService, jwtMidd)
	api.NewCompanywalletApi(app, companywalletService, jwtMidd)
	api.NewCompanyadminApi(app, companyadminService, jwtMidd)
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
	dbConnection.Close()
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
