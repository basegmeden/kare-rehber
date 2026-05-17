package models

import (
	"time"

	"gorm.io/gorm"
)

type PoolStatus string

const (
	PoolStatusPending  PoolStatus = "pending"
	PoolStatusApproved PoolStatus = "approved"
	PoolStatusRejected PoolStatus = "rejected"
)

type Coach struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Experience string         `json:"experience"`
	PoolStatus PoolStatus     `gorm:"default:pending" json:"pool_status"`
}
