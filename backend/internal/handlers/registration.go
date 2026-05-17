package handlers

import (
	"fmt"

	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/models"
	"kare-rehber/backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudent(c *fiber.Ctx) error {
	var body struct {
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Phone       string `json:"phone"`
		City        string `json:"city"`
		School      string `json:"school"`
		Grade       string `json:"grade"`
		DateOfBirth string `json:"date_of_birth"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if body.Name == "" || body.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and phone required"})
	}

	username := fmt.Sprintf("ogrenci_%s", body.Phone)
	tempPassword := body.Phone[len(body.Phone)-4:]
	hash, err := services.HashPassword(tempPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	user := models.User{
		Name:         body.Name,
		Surname:      body.Surname,
		Phone:        body.Phone,
		City:         body.City,
		Role:         models.RoleStudent,
		Username:     username,
		PasswordHash: hash,
		Status:       models.UserStatusPassive,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Phone number already registered"})
	}

	student := models.Student{
		UserID:             user.ID,
		RegistrationStatus: models.RegistrationPre,
		School:             body.School,
		Grade:              body.Grade,
	}
	database.DB.Create(&student)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ön kayıt alındı. Sizinle iletişime geçilecek.",
		"id":      user.ID,
	})
}

func RegisterCoach(c *fiber.Ctx) error {
	var body struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Phone      string `json:"phone"`
		City       string `json:"city"`
		Experience string `json:"experience"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if body.Name == "" || body.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and phone required"})
	}

	username := fmt.Sprintf("koc_%s", body.Phone)
	tempPassword := body.Phone[len(body.Phone)-4:]
	hash, err := services.HashPassword(tempPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}

	user := models.User{
		Name:         body.Name,
		Surname:      body.Surname,
		Phone:        body.Phone,
		City:         body.City,
		Role:         models.RoleCoach,
		Username:     username,
		PasswordHash: hash,
		Status:       models.UserStatusPassive,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Phone number already registered"})
	}

	coach := models.Coach{
		UserID:     user.ID,
		Experience: body.Experience,
		PoolStatus: models.PoolStatusPending,
	}
	database.DB.Create(&coach)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Başvurunuz alındı. İncelendikten sonra bilgilendirileceksiniz.",
		"id":      user.ID,
	})
}
