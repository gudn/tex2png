package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

type TemplateArgs struct {
  body string
}

func handler(c *fiber.Ctx) error {
	args := TemplateArgs {string(c.Body())}
	return c.SendString(fmt.Sprint(args))
}

func main() {
	app := fiber.New()
	app.Post("/", handler)
	app.Listen(fmt.Sprintf(":%s", os.Getenv("T2P_PORT")))
}
