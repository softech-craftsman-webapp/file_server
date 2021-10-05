package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Type      string         `gorm:"type:varchar(255);not null" json:"type"`
	Size      int64          `gorm:"type:bigint;not null" json:"size"`
	UserID    string         `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
