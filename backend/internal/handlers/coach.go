package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/models"
	"kare-rehber/backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type gormDatatypes = datatypes.JSON

// CoachListWeeks returns the active week + the one immediately before it.
func CoachListWeeks(c *fiber.Ctx) error {
	var weeks []models.Week
	database.DB.Order("week_number desc").Limit(2).Find(&weeks)
	return c.JSON(weeks)
}

func coachIDFromUser(userID uint) (uint, error) {
	var coach models.Coach
	if err := database.DB.Where("user_id = ?", userID).First(&coach).Error; err != nil {
		return 0, err
	}
	return coach.ID, nil
}

func CoachListStudents(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	coachID, err := coachIDFromUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coach profile not found"})
	}
	var assignments []models.CoachStudent
	database.DB.Preload("Student.User").Where("coach_id = ? AND is_active = true", coachID).Find(&assignments)
	return c.JSON(assignments)
}

func CoachListMeetings(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	coachID, err := coachIDFromUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coach profile not found"})
	}

	var meetings []models.Meeting
	database.DB.Preload("Student.User").Preload("Week").
		Where("coach_id = ?", coachID).
		Order("created_at desc").
		Find(&meetings)
	return c.JSON(meetings)
}

func CoachSubmitMeeting(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	coachID, err := coachIDFromUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coach profile not found"})
	}

	var body struct {
		StudentID  uint                     `json:"student_id"`
		WeekID     uint                     `json:"week_id"`
		Rating     int                      `json:"rating"`
		Notes      string                   `json:"notes"`
		Categories []map[string]interface{} `json:"categories"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Compute rating from categories average if provided, else use explicit rating.
	rating := body.Rating
	if len(body.Categories) > 0 {
		sum := 0
		for _, cat := range body.Categories {
			if r, ok := cat["rating"].(float64); ok {
				sum += int(r)
			}
		}
		avg := sum / len(body.Categories)
		if avg >= 1 && avg <= 5 {
			rating = avg
		}
	}
	if rating < 1 {
		rating = 1
	}
	if rating > 5 {
		rating = 5
	}

	if err := services.CanCoachEditMeeting(body.WeekID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	var count int64
	database.DB.Model(&models.CoachStudent{}).
		Where("coach_id = ? AND student_id = ? AND is_active = true", coachID, body.StudentID).
		Count(&count)
	if count == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Bu öğrenci size atanmamış"})
	}

	var detailsJSON gormDatatypes
	if len(body.Categories) > 0 {
		b, _ := json.Marshal(map[string]interface{}{"categories": body.Categories})
		detailsJSON = gormDatatypes(b)
	}

	// Upsert: update existing meeting for this coach+student+week if it exists
	var existing models.Meeting
	err = database.DB.Where("coach_id = ? AND student_id = ? AND week_id = ?", coachID, body.StudentID, body.WeekID).First(&existing).Error
	if err == nil {
		// Update existing
		updates := map[string]interface{}{
			"rating":     rating,
			"notes":      body.Notes,
			"details":    detailsJSON,
			"status":     models.MeetingStatusPending,
			"updated_at": time.Now(),
		}
		database.DB.Model(&existing).Updates(updates)
		return c.JSON(existing)
	}

	meeting := models.Meeting{
		StudentID: body.StudentID,
		CoachID:   coachID,
		WeekID:    body.WeekID,
		Rating:    rating,
		Notes:     body.Notes,
		Details:   detailsJSON,
		Status:    models.MeetingStatusPending,
	}
	if err := database.DB.Create(&meeting).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Görüşme kaydedilemedi"})
	}
	return c.Status(fiber.StatusCreated).JSON(meeting)
}

func CoachListMessages(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	var messages []models.Message
	database.DB.Preload("Sender").
		Where("sender_id = ? OR recipient_id = ?", userID, userID).
		Order("created_at desc").
		Find(&messages)
	return c.JSON(messages)
}

func CoachSendMessage(c *fiber.Ctx) error {
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

func CoachUpdateMeeting(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	coachID, err := coachIDFromUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Coach profile not found"})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	var meeting models.Meeting
	if err := database.DB.Where("id = ? AND coach_id = ?", id, coachID).First(&meeting).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Meeting not found"})
	}

	if err := services.CanCoachEditMeeting(meeting.WeekID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	var body struct {
		Rating int    `json:"rating"`
		Notes  string `json:"notes"`
	}
	c.BodyParser(&body)

	updates := map[string]interface{}{"status": models.MeetingStatusPending}
	if body.Rating >= 1 && body.Rating <= 5 {
		updates["rating"] = body.Rating
	}
	if body.Notes != "" {
		updates["notes"] = body.Notes
	}
	updates["updated_at"] = time.Now()

	database.DB.Model(&meeting).Updates(updates)
	return c.JSON(meeting)
}
