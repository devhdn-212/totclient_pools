package api

import (
	"context"
	"encoding/json"
	"gofibergocu/domain"
	"gofibergocu/dto"
	"gofibergocu/internal/connection"
	"gofibergocu/internal/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type customerApi struct {
	customerService domain.CustomerService
}

func NewCustomer(app *fiber.App,
	customerService domain.CustomerService,
	authmidle fiber.Handler) {
	ca := customerApi{
		customerService: customerService,
	}
	customers := app.Group("/customers", authmidle)

	customers.Get("", ca.Index)
	customers.Post("", ca.Create)
	customers.Put("/:id", ca.Update)
	customers.Delete("/:id", ca.Delete)
	customers.Get("/:id", ca.Show)
}

func (ca *customerApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ca.customerService.Index(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.JSON(dto.CreateResponseSuccess(res))
}
func (ca *customerApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Create Customer"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for CreateCustomer",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validation failed", fails))
	}
	err := ca.customerService.Create(c, req)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create customer",
			zap.String("code", req.Code),
			zap.String("name", req.Name),
			zap.String("message", err.Error()),
			zap.String("record", string(recordJson)),
		)
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusCreated).
		JSON(dto.CreateResponseSuccess(""))
}
func (ca *customerApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "UpdateCustomer"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for UpdateCustomer",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validation failed", fails))
	}

	// /customer/:id
	req.ID = ctx.Params("id")
	err := ca.customerService.Update(c, req)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to update customer",
			zap.String("id", req.ID),
			zap.String("error", err.Error()),
			zap.String("record", string(recordJson)),
		)
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	connection.Log.Info("Customer updated successfully",
		zap.String("id", req.ID),
	)
	return ctx.Status(http.StatusCreated).
		JSON(dto.CreateResponseSuccess(""))
}
func (ca *customerApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// /customer/:id
	id := ctx.Params("id")
	err := ca.customerService.Delete(c, id)
	if err != nil {
		connection.Log.Error("Failed to delete customer",
			zap.String("id", id),
			zap.String("error", err.Error()),
		)
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	connection.Log.Info("Customer deleted successfully",
		zap.String("id", id),
	)
	return ctx.SendStatus(http.StatusNoContent)
}
func (ca *customerApi) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// /customer/:id
	id := ctx.Params("id")
	data, err := ca.customerService.Show(c, id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(data))
}
