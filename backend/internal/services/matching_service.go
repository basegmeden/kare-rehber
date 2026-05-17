package services

import (
	"time"

	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/models"
)

// AssignCoachToStudents bulk-assigns a coach to a list of student IDs.
func AssignCoachToStudents(coachID uint, studentIDs []uint) error {
	now := time.Now()
	for _, sid := range studentIDs {
		// Deactivate any existing assignment first.
		database.DB.Model(&models.CoachStudent{}).
			Where("student_id = ? AND is_active = true", sid).
			Update("is_active", false)

		assignment := models.CoachStudent{
			CoachID:    coachID,
			StudentID:  sid,
			AssignedAt: now,
			IsActive:   true,
		}
		if err := database.DB.Create(&assignment).Error; err != nil {
			return err
		}
	}
	return nil
}

// AssignCoordinatorToStudents bulk-assigns a coordinator to a list of student IDs.
func AssignCoordinatorToStudents(coordinatorID uint, studentIDs []uint) error {
	now := time.Now()
	for _, sid := range studentIDs {
		database.DB.Model(&models.CoordinatorStudent{}).
			Where("student_id = ? AND is_active = true", sid).
			Update("is_active", false)

		assignment := models.CoordinatorStudent{
			CoordinatorID: coordinatorID,
			StudentID:     sid,
			AssignedAt:    now,
			IsActive:      true,
		}
		if err := database.DB.Create(&assignment).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetStudentsByCity returns confirmed students filtered by city.
func GetStudentsByCity(city string) ([]models.Student, error) {
	var students []models.Student
	q := database.DB.Preload("User").
		Joins("JOIN users ON users.id = students.user_id").
		Where("students.registration_status = ?", models.RegistrationConfirmed)
	if city != "" {
		q = q.Where("users.city = ?", city)
	}
	err := q.Find(&students).Error
	return students, err
}
