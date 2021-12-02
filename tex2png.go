package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"text/template"
)

var templates *template.Template

type TemplateArgs struct {
	Body string
}

func handler(c *fiber.Ctx) error {
	args := TemplateArgs{string(c.Body())}
	rendered := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(rendered, "tex", args); err == nil {
		return c.SendString(rendered.String())
	} else {
		log.Printf("templating error: %v", err)
		return c.SendStatus(500)
	}
}

func main() {
	var err error
	templates, err = template.ParseFiles("templates/tex")
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Post("/", handler)
	app.Listen(fmt.Sprintf(":%s", os.Getenv("T2P_PORT")))
}
