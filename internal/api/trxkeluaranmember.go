package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/dto"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/util"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type trxkeluaranmemberApi struct {
	trxkeluaranmemberService domain.TrxkeluaranmemberService
}

func NewTrxkeluaranmemberApi(app *fiber.App,
	trxkeluaranmemberService domain.TrxkeluaranmemberService,
	authmidle fiber.Handler) {
	ad := trxkeluaranmemberApi{
		trxkeluaranmemberService: trxkeluaranmemberService,
	}
	trxkeluaranmember := app.Group("/api/trxkeluaranmember", authmidle)
	trxkeluaranmember.Post("", ad.Index)
	trxkeluaranmember.Post("/save", ad.Save)
}
func (co *trxkeluaranmemberApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	datatoken := ctx.Locals("client_agen").(string)
	client_username := util.Parsing_final(datatoken)
	flag, _, idcomp := util.GetDataRedisClient(client_username)
	if !flag {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	flagpage := util.Validpage(client_username, "COMPANY-VIEW")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin - Trxkeluaranmember",
		})
	}

	var req dto.TrxkeluarandetailAll
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Create Admin"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for trxkeluaranmember",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}

	res, err := co.trxkeluaranmemberService.All(c, idcomp, req.IDtrxkeluaran)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"record":  res,
	})
}
func (co *trxkeluaranmemberApi) Save(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.TrxkeluaranmemberSave
	if err := ctx.BodyParser(&req); err != nil {
		connection.Log.Error("Failed to parse request body",
			zap.String("endpoint", "Create Trxkeluaran"),
			zap.String("body", string(ctx.Body())),
			zap.String("error", err.Error()),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)

	if len(fails) > 0 {
		connection.Log.Warn("Validation failed for update Trxkeluaranmember",
			zap.Any("validation_errors", fails),
			zap.Any("body", req),
		)
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData(http.StatusBadRequest, "validation failed", fails))
	}
	datatoken := ctx.Locals("client_agen").(string)
	client_username := util.Parsing_final(datatoken)
	flag, _, idcomp := util.GetDataRedisClient(client_username)
	fmt.Println("COMP : ", idcomp)
	if !flag {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, "internal server error"))
	}
	flagpage := util.Validpage(client_username, "COMPANY-SAVE")
	if !flagpage {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Please Contact Admin",
		})
	}

	err := co.trxkeluaranmemberService.Save(c, req, client_username, idcomp)
	if err != nil {
		recordJson, _ := json.Marshal(req)
		connection.Log.Error("Failed to create / update Trxkeluaranmember",
			zap.String("id", req.ID),
			zap.String("error", err.Error()),
			zap.String("record", string(recordJson)),
		)
		// cek duplicate entry
		if err.Error() == "duplicate entry" {
			return ctx.Status(http.StatusConflict).
				JSON(dto.CreateResponseError(http.StatusConflict, "Duplicate Entry"))
		}
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(http.StatusInternalServerError, err.Error()))
	}
	connection.Log.Info("Trxkeluaranmember create / update successfully",
		zap.String("id", req.ID),
	)

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))
}
