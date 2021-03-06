package main

import (
	"cowingram/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func main() {

	if !(strings.Compare(os.Getenv("PRODUCTION"), "true") == 0) {
		err := godotenv.Load()

		if err != nil {
			log.Panic("cowingram: error loading the environment")
		}
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/pin", database.CheckAvailability)

	if strings.Compare(os.Getenv("PRODUCTION"), "true") == 0 {
		err := app.Listen(":8080")
		if err != nil {
			return
		}
	} else {
		err := app.Listen(":3000")
		if err != nil {
			panic(err)
		}
	}
}
