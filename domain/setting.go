package domain

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

type Setting struct {
	ID               int            `db:"idsetting"`
	Appversion       string         `db:"appversion"`
	Startmaintenance pgtype.Time    `db:"startmaintenance"`
	Endmaintenance   pgtype.Time    `db:"endmaintenance"`
	Shio_parent      int            `db:"shio_parent"`
	Created          string         `db:"create_by"`
	CreatedAt        sql.NullTime   `db:"create_at"`
	Update           sql.NullString `db:"update_by"`
	UpdateAt         sql.NullTime   `db:"update_at"`
}

type SettingRepository interface {
	FindAll(ctx context.Context) ([]Setting, error)
	FindByID(ctx context.Context) (Setting, error)
}

// MaintenanceStatus is the maintenance window resolved against "now" —
// Start/End are formatted "HH:MM:SS" so callers (API responses) can display
// them directly without touching pgtype.
type MaintenanceStatus struct {
	Active bool
	Start  string
	End    string
}

type SettingService interface {
	CheckMaintenance(ctx context.Context) (MaintenanceStatus, error)
}
