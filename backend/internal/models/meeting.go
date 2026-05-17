package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MeetingStatus string

const (
	MeetingStatusPending  MeetingStatus = "pending"
	MeetingStatusApproved MeetingStatus = "approved"
	MeetingStatusRejected MeetingStatus = "rejected"
)

type Meeting struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	StudentID  uint           `gorm:"not null;index" json:"student_id"`
	Student    Student        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	CoachID    uint           `gorm:"not null;index" json:"coach_id"`
	Coach      Coach          `gorm:"foreignKey:CoachID" json:"coach,omitempty"`
	WeekID     uint           `gorm:"not null;index" json:"week_id"`
	Week       Week           `gorm:"foreignKey:WeekID" json:"week,omitempty"`
	Rating     int            `gorm:"check:rating >= 1 AND rating <= 5" json:"rating"`
	Notes      string         `gorm:"type:text" json:"notes"`
	Details    datatypes.JSON `gorm:"type:json" json:"details,omitempty"`
	Status     MeetingStatus  `gorm:"default:pending" json:"status"`
	ApprovedBy *uint          `json:"approved_by,omitempty"`
	ApprovedAt *time.Time     `json:"approved_at,omitempty"`
}

// MeetingLog stores old data when admin edits a meeting.
type MeetingLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	MeetingID uint           `gorm:"not null;index" json:"meeting_id"`
	ChangedBy uint           `gorm:"not null" json:"changed_by"`
	OldData   datatypes.JSON `gorm:"type:jsonb" json:"old_data"`
	NewData   datatypes.JSON `gorm:"type:jsonb" json:"new_data"`
}
