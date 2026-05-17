package handlers

import (
	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ParentListMeetings(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)

	var student models.Student
	if err := database.DB.Where("parent_user_id = ?", userID).First(&student).Error; err != nil {
		return c.JSON([]models.Meeting{})
	}

	var meetings []models.Meeting
	database.DB.Preload("Coach.User").Preload("Week").Preload("Student.User").
		Where("student_id = ? AND status = ?", student.ID, models.MeetingStatusApproved).
		Order("created_at desc").
		Find(&meetings)
	return c.JSON(meetings)
}

func ParentListMessages(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	var messages []models.Message
	database.DB.Preload("Sender").
		Where("sender_id = ? OR recipient_id = ?", userID, userID).
		Order("created_at desc").
		Find(&messages)
	return c.JSON(messages)
}

func ParentSendMessage(c *fiber.Ctx) error {
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
