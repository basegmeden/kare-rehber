package services

import (
	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/models"
)

type OverviewReport struct {
	TotalStudents     int64 `json:"total_students"`
	ActiveStudents    int64 `json:"active_students"`
	PassiveStudents   int64 `json:"passive_students"`
	TotalCoaches      int64 `json:"total_coaches"`
	TotalMeetings     int64 `json:"total_meetings"`
	PendingMeetings   int64 `json:"pending_meetings"`
	ApprovedMeetings  int64 `json:"approved_meetings"`
}

type CityReport struct {
	City          string `json:"city"`
	StudentCount  int64  `json:"student_count"`
	MeetingCount  int64  `json:"meeting_count"`
}

type MentorPerformance struct {
	CoachID     uint   `json:"coach_id"`
	CoachName   string `json:"coach_name"`
	StudentCount int64 `json:"student_count"`
	MeetingCount int64 `json:"meeting_count"`
	MissingCount int64 `json:"missing_count"`
}

func GetOverviewReport() (*OverviewReport, error) {
	var r OverviewReport
	database.DB.Model(&models.Student{}).Count(&r.TotalStudents)
	database.DB.Model(&models.Student{}).
		Joins("JOIN users ON users.id = students.user_id").
		Where("users.status = ?", models.UserStatusActive).
		Count(&r.ActiveStudents)
	r.PassiveStudents = r.TotalStudents - r.ActiveStudents
	database.DB.Model(&models.Coach{}).Where("pool_status = ?", models.PoolStatusApproved).Count(&r.TotalCoaches)
	database.DB.Model(&models.Meeting{}).Count(&r.TotalMeetings)
	database.DB.Model(&models.Meeting{}).Where("status = ?", models.MeetingStatusPending).Count(&r.PendingMeetings)
	database.DB.Model(&models.Meeting{}).Where("status = ?", models.MeetingStatusApproved).Count(&r.ApprovedMeetings)
	return &r, nil
}

func GetCityReport() ([]CityReport, error) {
	var results []CityReport
	err := database.DB.
		Model(&models.Student{}).
		Select("users.city, COUNT(DISTINCT students.id) as student_count, COUNT(meetings.id) as meeting_count").
		Joins("JOIN users ON users.id = students.user_id").
		Joins("LEFT JOIN coach_students cs ON cs.student_id = students.id AND cs.is_active = true").
		Joins("LEFT JOIN meetings ON meetings.student_id = students.id").
		Group("users.city").
		Scan(&results).Error
	return results, err
}

func GetMissingMeetingsReport() ([]models.Coach, error) {
	return GetCoachAlerts()
}

func GetMentorPerformanceReport() ([]MentorPerformance, error) {
	var results []MentorPerformance
	err := database.DB.
		Model(&models.Coach{}).
		Select(`coaches.id as coach_id,
			(users.name || ' ' || users.surname) as coach_name,
			COUNT(DISTINCT cs.student_id) as student_count,
			COUNT(DISTINCT m.id) as meeting_count`).
		Joins("JOIN users ON users.id = coaches.user_id").
		Joins("LEFT JOIN coach_students cs ON cs.coach_id = coaches.id AND cs.is_active = true").
		Joins("LEFT JOIN meetings m ON m.coach_id = coaches.id").
		Group("coaches.id, users.name, users.surname").
		Scan(&results).Error

	for i := range results {
		results[i].MissingCount = results[i].StudentCount - results[i].MeetingCount
		if results[i].MissingCount < 0 {
			results[i].MissingCount = 0
		}
	}
	return results, err
}
