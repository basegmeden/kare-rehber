package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	SenderID    uint           `gorm:"not null;index" json:"sender_id"`
	Sender      User           `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	RecipientID uint           `gorm:"not null;index" json:"recipient_id"`
	Recipient   User           `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
	Body        string         `gorm:"type:text;not null" json:"body"`
	IsRead      bool           `gorm:"default:false" json:"is_read"`
}
