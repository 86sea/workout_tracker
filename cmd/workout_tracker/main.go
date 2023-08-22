package main

import (
	"github.com/gofiber/fiber/v2"
	"workout_tracker/internal/database"
    "workout_tracker/internal/routes"
    //"log"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
    database.Connect()
    routes.AuthRoutes(app)
    routes.WebRoutes(app)
    app.Listen(":3000")
}
