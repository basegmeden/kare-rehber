package models

import (
	"time"

	"gorm.io/gorm"
)

type RegistrationStatus string

const (
	RegistrationPre       RegistrationStatus = "pre"
	RegistrationConfirmed RegistrationStatus = "confirmed"
)

type Student struct {
	ID                 uint               `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`
	UserID             uint               `gorm:"uniqueIndex;not null" json:"user_id"`
	User               User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentUserID       *uint              `json:"parent_user_id,omitempty"`
	RegistrationStatus RegistrationStatus `gorm:"default:pre" json:"registration_status"`
	School             string             `json:"school"`
	Grade              string             `json:"grade"`
	Notes              string             `json:"notes"`
}

type CoachStudent struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CoachID    uint      `gorm:"not null;index" json:"coach_id"`
	Coach      Coach     `gorm:"foreignKey:CoachID" json:"coach,omitempty"`
	StudentID  uint      `gorm:"not null;index" json:"student_id"`
	Student    Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
}

type CoordinatorStudent struct {
	ID            uint        `gorm:"primarykey" json:"id"`
	CoordinatorID uint        `gorm:"not null;index" json:"coordinator_id"`
	Coordinator   Coordinator `gorm:"foreignKey:CoordinatorID" json:"coordinator,omitempty"`
	StudentID     uint        `gorm:"not null;index" json:"student_id"`
	Student       Student     `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	AssignedAt    time.Time   `json:"assigned_at"`
	IsActive      bool        `gorm:"default:true" json:"is_active"`
}
