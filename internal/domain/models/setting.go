package models

type SystemSetting struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Key         string `gorm:"uniqueIndex;not null" json:"key"` // e.g. "maintenance_mode"
	Value       string `json:"value"`                           // "true" or "false"
	Description string `json:"description"`
}