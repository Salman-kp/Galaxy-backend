package models

type RoleWage struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Role string `gorm:"size:50;uniqueIndex;not null" json:"role"`
	Wage int64  `gorm:"not null" json:"wage"`
}