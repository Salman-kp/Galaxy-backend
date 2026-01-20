package models

type Permission struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Slug        string `gorm:"size:100;uniqueIndex;not null" json:"slug"` 
    Description string `gorm:"size:255" json:"description"`
}

type AdminRole struct {
    ID          uint         `gorm:"primaryKey" json:"id"`
    Name        string       `gorm:"size:100;uniqueIndex;not null" json:"name"`
    Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"permissions"`
}