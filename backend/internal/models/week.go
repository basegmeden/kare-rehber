package models

import (
	"time"

	"gorm.io/gorm"
)

type Week struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	WeekNumber int            `gorm:"not null" json:"week_number"`
	Label      string         `gorm:"not null" json:"label"` // "3. Hafta (11-17 Mayıs)"
	StartDate  time.Time      `gorm:"not null" json:"start_date"`
	EndDate    time.Time      `gorm:"not null" json:"end_date"`
	IsActive   bool           `gorm:"default:false" json:"is_active"`
	IsLocked   bool           `gorm:"default:false" json:"is_locked"`
	CreatedBy  uint           `json:"created_by"`
}
