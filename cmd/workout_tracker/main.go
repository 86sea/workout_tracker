package main

import (
	"github.com/gofiber/fiber/v2"
	"workout_tracker/internal/database"
    //"log"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
    database.Connect()
	//if err := database.Connect(); err != nil {
		//log.Fatal(err)
	//} else{
        //log.Println("connected to db")
    //}
    app.Listen(":3000")
}
