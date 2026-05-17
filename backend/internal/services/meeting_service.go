package services

import (
	"errors"
	"time"

	"encoding/json"

	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/models"

	"gorm.io/datatypes"
)

var ErrWeekLocked = errors.New("bu hafta kilitli, düzenleme yapılamaz")
var ErrEditWindowExpired = errors.New("yalnızca son 1 hafta öncesine kadar düzenleme yapılabilir")

// CanCoachEditMeeting checks the 1-week-prior edit limit for coaches.
func CanCoachEditMeeting(weekID uint) error {
	var week models.Week
	if err := database.DB.First(&week, weekID).Error; err != nil {
		return err
	}
	if week.IsLocked {
		return ErrWeekLocked
	}
	// Allow editing current week and one prior week based on end dates.
	cutoff := time.Now().AddDate(0, 0, -14)
	if week.EndDate.Before(cutoff) {
		return ErrEditWindowExpired
	}
	return nil
}

// AdminEditMeeting updates a meeting and writes old data to meeting_logs.
func AdminEditMeeting(meetingID, adminID uint, rating int, notes string) error {
	var meeting models.Meeting
	if err := database.DB.First(&meeting, meetingID).Error; err != nil {
		return err
	}

	oldRaw, _ := json.Marshal(map[string]interface{}{
		"rating": meeting.Rating,
		"notes":  meeting.Notes,
		"status": string(meeting.Status),
	})
	newRaw, _ := json.Marshal(map[string]interface{}{
		"rating": rating,
		"notes":  notes,
		"status": string(meeting.Status),
	})

	logEntry := models.MeetingLog{
		MeetingID: meetingID,
		ChangedBy: adminID,
		OldData:   datatypes.JSON(oldRaw),
		NewData:   datatypes.JSON(newRaw),
	}
	if err := database.DB.Create(&logEntry).Error; err != nil {
		return err
	}

	return database.DB.Model(&meeting).Updates(map[string]interface{}{
		"rating": rating,
		"notes":  notes,
	}).Error
}

// GetCoachAlerts returns coaches who have students without meetings for the active week.
func GetCoachAlerts() ([]models.Coach, error) {
	var activeWeek models.Week
	if err := database.DB.Where("is_active = true").First(&activeWeek).Error; err != nil {
		return nil, err
	}

	var coaches []models.Coach
	subquery := database.DB.
		Model(&models.Meeting{}).
		Select("coach_id").
		Where("week_id = ?", activeWeek.ID)

	err := database.DB.
		Preload("User").
		Joins("JOIN coach_students cs ON cs.coach_id = coaches.id AND cs.is_active = true").
		Where("coaches.id NOT IN (?)", subquery).
		Group("coaches.id").
		Find(&coaches).Error

	return coaches, err
}
