package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("tex2png")
	})
	app.Listen(fmt.Sprintf(":%s", os.Getenv("T2P_PORT")))
}
