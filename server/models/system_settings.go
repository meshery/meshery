package models

import (
	"encoding/json"
	"time"

	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
	"gorm.io/gorm/clause"
)

// SystemSetting is a server-wide key/value setting persisted in Meshery
// Server's own database. Settings scoped here describe how this Meshery
// Server deployment behaves (for example, the default Meshery Operator /
// MeshSync / Broker configuration applied to every managed cluster) and are
// deliberately independent of the provider: they belong to the server
// instance, not to a user or an organization.
type SystemSetting struct {
	Key       string    `json:"key" gorm:"primaryKey"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName overrides the GORM default so the table reads naturally.
func (SystemSetting) TableName() string {
	return "system_settings"
}

// ControllersConfigDefaultsSettingKey is the system_settings key under which
// the server-wide MesheryControllersConfig defaults are stored.
const ControllersConfigDefaultsSettingKey = "controllers_config_defaults"

// GetSystemSetting reads a raw setting value. The second return reports
// whether the key exists.
func GetSystemSetting(dbHandler *database.Handler, key string) (string, bool, error) {
	if dbHandler == nil {
		return "", false, ErrSystemSettings(ErrDBConnection)
	}
	dbHandler.Lock()
	defer dbHandler.Unlock()

	var settings []SystemSetting
	result := dbHandler.Where("key = ?", key).Limit(1).Find(&settings)
	if result.Error != nil {
		return "", false, ErrSystemSettings(result.Error)
	}
	if len(settings) == 0 {
		return "", false, nil
	}
	return settings[0].Value, true, nil
}

// SetSystemSetting upserts a raw setting value.
func SetSystemSetting(dbHandler *database.Handler, key, value string) error {
	if dbHandler == nil {
		return ErrSystemSettings(ErrDBConnection)
	}
	dbHandler.Lock()
	defer dbHandler.Unlock()

	setting := SystemSetting{Key: key, Value: value, UpdatedAt: time.Now().UTC()}
	result := dbHandler.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&setting)
	if result.Error != nil {
		return ErrSystemSettings(result.Error)
	}
	return nil
}

// GetControllersConfigDefaults returns the persisted server-wide controllers
// configuration defaults, or nil when none have been set.
func GetControllersConfigDefaults(dbHandler *database.Handler) (*controllersconfig.MesheryControllersConfig, error) {
	value, exists, err := GetSystemSetting(dbHandler, ControllersConfigDefaultsSettingKey)
	if err != nil {
		return nil, err
	}
	if !exists || value == "" {
		return nil, nil
	}
	cfg := &controllersconfig.MesheryControllersConfig{}
	if err := json.Unmarshal([]byte(value), cfg); err != nil {
		return nil, ErrSystemSettings(err)
	}
	return cfg, nil
}

// SaveControllersConfigDefaults persists the server-wide controllers
// configuration defaults. A nil or empty document clears the stored defaults.
func SaveControllersConfigDefaults(dbHandler *database.Handler, cfg *controllersconfig.MesheryControllersConfig) error {
	if cfg == nil || (cfg.Operator == nil && cfg.Meshsync == nil && cfg.Broker == nil) {
		return SetSystemSetting(dbHandler, ControllersConfigDefaultsSettingKey, "")
	}
	cfg.SchemaVersion = connections.ControllersConfigSchemaVersion
	encoded, err := json.Marshal(cfg)
	if err != nil {
		return ErrSystemSettings(err)
	}
	return SetSystemSetting(dbHandler, ControllersConfigDefaultsSettingKey, string(encoded))
}
