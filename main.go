package main

import (
	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gofibergocu/dto"
	"gofibergocu/internal/api"
	"gofibergocu/internal/config"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/repository"
	"gofibergocu/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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
	app.Use(logger.New())

	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/"
		},
		Format: "${time} | ${status} | ${latency} | ${ips[0]} | ${method} | ${path} - ${queryParams} ${body}\n",
	}))
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
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	log.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	log.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	dbConnection.Close()
	// redisConn.Close()
	log.Println("Fiber was successful shutdown.")
}
