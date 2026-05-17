package handlers

import (
	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Students cannot view their own meeting notes — only messages.

func StudentGetProfile(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)

	var student models.Student
	if err := database.DB.Preload("User").Where("user_id = ?", userID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	var coachStudent models.CoachStudent
	coachInfo := fiber.Map{}
	if err := database.DB.Preload("Coach.User").Where("student_id = ? AND is_active = ?", student.ID, true).First(&coachStudent).Error; err == nil {
		coachInfo = fiber.Map{
			"name":    coachStudent.Coach.User.Name + " " + coachStudent.Coach.User.Surname,
			"user_id": coachStudent.Coach.UserID,
		}
	}

	return c.JSON(fiber.Map{
		"id":                  student.ID,
		"name":                student.User.Name,
		"surname":             student.User.Surname,
		"city":                student.User.City,
		"school":              student.School,
		"grade":               student.Grade,
		"registration_status": student.RegistrationStatus,
		"coach":               coachInfo,
	})
}

func StudentListMessages(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	var messages []models.Message
	database.DB.Preload("Sender").
		Where("sender_id = ? OR recipient_id = ?", userID, userID).
		Order("created_at desc").
		Find(&messages)
	return c.JSON(messages)
}

func StudentSendMessage(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	var body struct {
		RecipientID uint   `json:"recipient_id"`
		Body        string `json:"body"`
	}
	if err := c.BodyParser(&body); err != nil || body.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "recipient_id and body required"})
	}
	msg := models.Message{
		SenderID:    userID,
		RecipientID: body.RecipientID,
		Body:        body.Body,
	}
	database.DB.Create(&msg)
	return c.Status(fiber.StatusCreated).JSON(msg)
}
