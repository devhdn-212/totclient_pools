package dto

type SettingData struct {
	ID               int    `json:"setting_id"`
	Appversion       string `json:"setting_appversion"`
	Startmaintenance string `json:"setting_startmaintenance"`
	Endmaintenance   string `json:"setting_endmaintenance"`
	Shio_parent      int    `json:"setting_shio_parent"`
}
