package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
)

var templates *template.Template

type TemplateArgs struct {
	Body string
}

func tex2png(texCode *bytes.Buffer) (string, error) {
	tmpDir, err := os.MkdirTemp("", "t2p-")
	if err != nil {
		return "", err
	}
	srcFn := path.Join(tmpDir, "source.tex")
	os.WriteFile(srcFn, texCode.Bytes(), 0o664)
	cmd := exec.Command("pdflatex", "-halt-on-error", srcFn)
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}
	cmd = exec.Command("pdftoppm", path.Join(tmpDir, "source.pdf"), "output", "-png")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return path.Join(tmpDir, "output-1.png"), nil
}

func handler(c *fiber.Ctx) error {
	args := TemplateArgs{string(c.Body())}
	rendered := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(rendered, "tex", args); err == nil {
		png, err := tex2png(rendered)
		if err != nil {
			log.Printf("rendering error: %v", err)
			return c.SendStatus(500)
		}
		defer os.RemoveAll(path.Dir(png))
		return c.SendFile(png)
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
