package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devhdn-212/totclient_pools/domain"
	"github.com/devhdn-212/totclient_pools/internal/connection"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// respondIfMaintenance checks the maintenance window and, if active, writes
// the 503 response itself and reports handled=true so the caller returns
// immediately without running its normal business logic. A failure to check
// the status (e.g. Redis/DB hiccup) fails open — it's logged but does not
// block the request, since a broken settings lookup shouldn't take down
// login/checkout entirely.
func respondIfMaintenance(ctx *fiber.Ctx, c context.Context, settingService domain.SettingService, endpoint string) (handled bool, err error) {
	maintenance, err := settingService.CheckMaintenance(c)
	if err != nil {
		connection.Log.Error("Failed to check maintenance status",
			zap.String("endpoint", endpoint),
			zap.String("error", err.Error()),
		)
		return false, nil
	}
	if !maintenance.Active {
		return false, nil
	}

	writeErr := ctx.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
		"status":            fiber.StatusServiceUnavailable,
		"message":           fmt.Sprintf("Website sedang maintenance dari jam %s sampai %s", maintenance.Start, maintenance.End),
		"maintenance":       true,
		"start_maintenance": maintenance.Start,
		"end_maintenance":   maintenance.End,
	})
	return true, writeErr
}
