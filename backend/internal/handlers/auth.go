package handlers

import (
	"kare-rehber/backend/internal/middleware"
	"kare-rehber/backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username and password required"})
	}

	user, err := services.FindUserByUsername(body.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error"})
	}
	if user == nil || !services.CheckPassword(user.PasswordHash, body.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token generation failed"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"name":     user.Name,
			"surname":  user.Surname,
			"username": user.Username,
			"role":     user.Role,
			"city":     user.City,
		},
	})
}

func Me(c *fiber.Ctx) error {
	userID := middleware.CurrentUserID(c)
	user, err := services.FindUserByID(userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"surname":  user.Surname,
		"username": user.Username,
		"role":     user.Role,
		"city":     user.City,
		"phone":    user.Phone,
		"status":   user.Status,
	})
}
