package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devhdn-212/totclient_api/domain"
	"github.com/devhdn-212/totclient_api/internal/connection"
	"github.com/devhdn-212/totclient_api/internal/util"
)

const (
	RedisSettingKey = "client:setting"
	settingCacheTTL = 5 * time.Minute
)

// settingCache is the small JSON-safe shape actually cached in Redis —
// pgtype.Time doesn't round-trip through encoding/json reliably, so the
// maintenance window is stored pre-formatted as "HH:MM:SS" strings instead
// of caching domain.Setting directly.
type settingCache struct {
	Appversion string `json:"appversion"`
	Start      string `json:"start"`
	End        string `json:"end"`
	ShioParent int    `json:"shio_parent"`
}

type settingService struct {
	repo domain.SettingRepository
}

func NewSettingService(repo domain.SettingRepository) domain.SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) getSetting(ctx context.Context) (settingCache, error) {
	cached, found, err := connection.GetRedis(RedisSettingKey)
	if err == nil && found {
		var data settingCache
		if jsonErr := json.Unmarshal([]byte(cached), &data); jsonErr == nil {
			return data, nil
		}
	}

	setting, err := s.repo.FindByID(ctx)
	if err != nil {
		return settingCache{}, err
	}

	data := settingCache{
		Appversion: setting.Appversion,
		Start:      util.PgTimeToString(setting.Startmaintenance),
		End:        util.PgTimeToString(setting.Endmaintenance),
		ShioParent: setting.Shio_parent,
	}
	go connection.SetRedis(RedisSettingKey, data, settingCacheTTL)
	return data, nil
}

// CheckMaintenance compares "now" (Jakarta time) against today's date
// combined with the configured Startmaintenance/Endmaintenance time-of-day.
// An end time earlier than the start time is treated as crossing midnight
// (maintenance window runs into the next day).
func (s *settingService) CheckMaintenance(ctx context.Context) (domain.MaintenanceStatus, error) {
	data, err := s.getSetting(ctx)
	if err != nil {
		return domain.MaintenanceStatus{}, err
	}

	if data.Start == "" || data.End == "" {
		return domain.MaintenanceStatus{}, nil
	}

	now := util.GetNowJakarta()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, util.LocJakarta)

	start, err := time.ParseInLocation("15:04:05", data.Start, util.LocJakarta)
	if err != nil {
		return domain.MaintenanceStatus{}, nil
	}
	end, err := time.ParseInLocation("15:04:05", data.End, util.LocJakarta)
	if err != nil {
		return domain.MaintenanceStatus{}, nil
	}

	startAt := today.Add(time.Duration(start.Hour())*time.Hour + time.Duration(start.Minute())*time.Minute + time.Duration(start.Second())*time.Second)
	endAt := today.Add(time.Duration(end.Hour())*time.Hour + time.Duration(end.Minute())*time.Minute + time.Duration(end.Second())*time.Second)

	if endAt.Before(startAt) {
		endAt = endAt.Add(24 * time.Hour)
	}

	active := !now.Before(startAt) && now.Before(endAt)

	// Full date+time (not just "HH:MM") — endAt may land on the next
	// calendar day when the window crosses midnight, so the date matters.
	const displayLayout = "2006-01-02 15:04"
	return domain.MaintenanceStatus{
		Active: active,
		Start:  startAt.Format(displayLayout),
		End:    endAt.Format(displayLayout),
	}, nil
}
