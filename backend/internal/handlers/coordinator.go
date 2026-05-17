package handlers

import (
	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func coordinatorIDFromUser(userID uint) (uint, error) {
	var coord models.Coordinator
	if err := database.DB.Where("user_id = ?", userID).First(&coord).Error; err != nil {
		return 0, err
	}
	return coord.ID, nil
}

func CoordinatorListStudents(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	coordID, err := coordinatorIDFromUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coordinator profile not found"})
	}
	var assignments []models.CoordinatorStudent
	database.DB.Preload("Student.User").Where("coordinator_id = ? AND is_active = true", coordID).Find(&assignments)
	return c.JSON(assignments)
}

func CoordinatorListMeetings(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	coordID, err := coordinatorIDFromUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coordinator profile not found"})
	}

	// Get student IDs for this coordinator.
	var assignments []models.CoordinatorStudent
	database.DB.Where("coordinator_id = ? AND is_active = true", coordID).Find(&assignments)
	studentIDs := make([]uint, len(assignments))
	for i, a := range assignments {
		studentIDs[i] = a.StudentID
	}

	var meetings []models.Meeting
	if len(studentIDs) > 0 {
		database.DB.Preload("Student.User").Preload("Coach.User").Preload("Week").
			Where("student_id IN (?)", studentIDs).
			Order("created_at desc").
			Find(&meetings)
	}
	return c.JSON(meetings)
}
