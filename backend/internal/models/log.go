package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)


type SMSType string

const (
	SMSTypeCredentials SMSType = "credentials"
	SMSTypeBulk        SMSType = "bulk"
	SMSTypeAlert       SMSType = "alert"
	SMSTypeCustom      SMSType = "custom"
)

type SMSStatus string

const (
	SMSStatusSent   SMSStatus = "sent"
	SMSStatusFailed SMSStatus = "failed"
)

type SMSLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	RecipientID *uint     `json:"recipient_id,omitempty"`
	Phone       string    `gorm:"not null" json:"phone"`
	Message     string    `gorm:"type:text;not null" json:"message"`
	Type        SMSType   `json:"type"`
	Status      SMSStatus `json:"status"`
}

type AuditLog struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     string         `gorm:"not null" json:"action"`
	EntityType string         `json:"entity_type"`
	EntityID   *uint          `json:"entity_id,omitempty"`
	Details    datatypes.JSON `gorm:"type:jsonb" json:"details"`
}
