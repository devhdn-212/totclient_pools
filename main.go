package main

import (
	"encoding/json"
	"gofibergocu/dto"
	"gofibergocu/internal/api"
	"gofibergocu/internal/config"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/repository"
	"gofibergocu/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		fields := []zap.Field{
			zap.String("request_id", c.Locals("requestid").(string)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency_ms", latency),
			zap.String("ip", c.IP()),
		}
		if c.Method() == fiber.MethodPost && c.Is("json") {
			body := c.Body()
			if len(body) > 0 {
				// bisa simpan sebagai string
				fields = append(fields, zap.String("body", string(body)))

				// atau parse JSON agar structured
				var parsed map[string]interface{}
				if err := json.Unmarshal(body, &parsed); err == nil {
					fields = append(fields, zap.Any("json_body", parsed))
				}
			}
		}

		// log
		logger.Info("request", fields...)
		return err
	})
	app.Use(requestid.New())
	app.Use(etag.New())

	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).
				JSON(dto.CreateResponseError("missing token, please login"))
		},
	})
	goquExec := repository.NewGoquExecutor(dbConnection)
	customerRepository := repository.NewCustomer(goquExec)
	userRepository := repository.NewUser(dbConnection)
	customerService := service.NewCustomerService(dbConnection, customerRepository)
	authService := service.NewAuth(cnf, userRepository)

	api.NewCustomer(app, customerService, jwtMidd)
	api.NewAuth(app, authService)

	go func() {
		appsPort := cnf.Server.Port
		if appsPort == "" {
			appsPort = "7072"
		}

		if err := app.Listen(":" + appsPort); err != nil {

			logger.Fatal("Failed to start app", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	logger.Info("Gracefully shutting down...")
	_ = app.Shutdown()

	logger.Info("Running cleanup tasks...")

	// Your cleanup tasks go here
	dbConnection.Close()
	connection.RDB.Close()
	logger.Info("Fiber was successful shutdown.")
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
