package middleware

import (
	"github.com/gofiber/fiber/v2"
	"workout_tracker/internal/jwt"
)

func RequireAuth(c *fiber.Ctx) error {
	t := c.Cookies("jwt")
	_, err := jwt.Verify(t)
	if err != nil {
		return c.Redirect("/login")
	}

	// If the JWT is valid, continue with the request
	return c.Next()
}
