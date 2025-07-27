package gorm

import (
	"time"
)

type GormUser struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	Username       string    `gorm:"uniqueIndex;not null"`
	Email          string    `gorm:"uniqueIndex;not null"`
	HashedPassword string    `gorm:"not null"`
	RoleID         uint      `gorm:"not null"`
	Role           GormRole  `gorm:"foreignKey:RoleID;references:ID"`
	LastSeen       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
