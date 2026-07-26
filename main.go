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
	"github.com/devhdn-212/totclient_api/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cnf := config.Get()
	logger := NewGCPLogger(cnf.Telegram)

	connection.SetLogger(logger)

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
	// Fiber v2 does not recover from panics on its own — without this, a
	// panic anywhere in a handler takes down the whole process instead of
	// just failing that one request. Routed through connection.Log.Error so
	// it also fires the Telegram alert hook below.
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			connection.Log.Error("panic recovered",
				zap.Any("panic", e),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		},
	}))
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
	riwayatTransaksiService := service.NewRiwayatTransaksiService(pasaranService, trxkeluarandetailService, trxkeluaranmemberService, memberinfoService)
	resultService := service.NewResultService(pasaranService, trxkeluaranRepository, memberinfoService)

	api.NewMemberInfo(app, memberinfoService, settingService)
	api.NewCheckoutApi(app, checkoutService, settingService)
	api.NewRiwayatTransaksiApi(app, riwayatTransaksiService)
	api.NewResultApi(app, resultService)
	api.NewServiceinit(app, pasaranService)

	go func() {
		appsPort := cnf.Server.Port
		if err := app.Listen(":" + appsPort); err != nil {
			logger.Fatal("Failed to start app", zap.Error(err))
		}
	}()

	// Drains periods checkout marked dirty (service.MarkTotalsDirty) and
	// recomputes their total_member/total_bet/total_pairs/total_payout — see
	// trxkeluaran.go for why this replaced an inline per-checkout increment
	// (that serialized every concurrent checkout for a period on one row
	// lock). These numbers only feed a display-only agen dashboard (not live
	// business logic — limit checks use their own separate Redis counters),
	// so 1 minute of staleness is fine; a longer interval also means more
	// checkouts for a busy period get batched into a single recompute.
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			service.FlushDirtyTotals(context.Background(), trxkeluaranRepository)
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

func NewGCPLogger(tg config.Telegram) *zap.Logger {
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

	return zap.New(core, zap.AddCaller(), zap.Hooks(func(entry zapcore.Entry) error {
		// Every Error/Fatal log anywhere in the app (including panics
		// recovered above) flows through here — entry.Caller pinpoints the
		// exact file:line so the alert says where to look, not just what
		// broke. Fire-and-forget: never let a slow Telegram API delay the
		// request that produced the log.
		if entry.Level < zapcore.ErrorLevel {
			return nil
		}
		go util.SendTelegramAlert(tg, entry.Level.CapitalString(), entry.Caller.String(), entry.Message)
		return nil
	}))
}
