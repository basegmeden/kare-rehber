package middleware

import (
	"kare-rehber/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := models.Role(c.Locals("role").(string))
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}
}

func CurrentUserID(c *fiber.Ctx) uint {
	return c.Locals("userID").(uint)
}

func CurrentRole(c *fiber.Ctx) models.Role {
	return models.Role(c.Locals("role").(string))
}
