package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/api"
	"github.com/devhdn-212/totclient_api/internal/config"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/repository"
	"github.com/devhdn-212/totclient_api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
	defer connection.RDBLimit.Close()

	if !connection.RedisHealth() {
		logger.Fatal("Redis is not healthy")
	}

	app := fiber.New()
	app.Use(requestid.New())
	app.Use(etag.New())
	app.Use(limiter.New(limiter.Config{
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

	pgxExec := repository.NewPGXExecutor(dbPool)

	pasaranRepository := repository.NewPasaranRepository(pgxExec)
	trxkeluaranRepository := repository.NewTrxkeluaranRepository(pgxExec)
	trxkeluarandetailRepository := repository.NewTrxkeluarandetailRepository(pgxExec)
	trxkeluaranmemberRepository := repository.NewTrxkeluaranmemberRepository(pgxExec)
	settingRepository := repository.NewSettingRepository(pgxExec)
	pasaranService := service.NewPasaranService(dbPool, pasaranRepository, trxkeluaranRepository)
	memberinfoService := service.NewMemberinfoService()
	settingService := service.NewSettingService(settingRepository)

	checkoutService := service.NewCheckoutService(dbPool, trxkeluaranRepository, pasaranService, memberinfoService)
	trxkeluarandetailService := service.NewTrxkeluarandetailService(dbPool, trxkeluarandetailRepository)
	trxkeluaranmemberService := service.NewTrxkeluaranmemberService(dbPool, trxkeluaranmemberRepository)
	riwayatTransaksiService := service.NewRiwayatTransaksiService(pasaranService, trxkeluaranRepository, trxkeluarandetailService, trxkeluaranmemberService, memberinfoService)

	api.NewMemberInfo(app, memberinfoService, settingService)
	api.NewCheckoutApi(app, checkoutService, settingService)
	api.NewRiwayatTransaksiApi(app, riwayatTransaksiService)
	api.NewServiceinit(app, pasaranService)

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
