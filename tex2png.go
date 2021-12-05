package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"syscall"
	"text/template"
)

var templates *template.Template
var inputRegexp *regexp.Regexp

type TemplateArgs struct {
	Body string
}

func testInput(texCode *bytes.Buffer) bool {
	return inputRegexp.Find(texCode.Bytes()) == nil
}

func tex2png(texCode *bytes.Buffer) (string, error) {
	tmpDir, err := os.MkdirTemp("", "t2p-")
	if err != nil {
		return "", err
	}
	srcFn := path.Join(tmpDir, "source.tex")
	os.WriteFile(srcFn, texCode.Bytes(), 0o664)
	cmd := exec.Command("pdflatex", "-halt-on-error", "-no-shell-escape", srcFn)
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
	inputRegexp = regexp.MustCompile("\\\\(input|include|openin|openout)")
	rendered := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(rendered, "tex", args); err == nil {
		if !testInput(rendered) {
			c.SendString("Illegal input")
			return c.SendStatus(400)
		}
		png, err := tex2png(rendered)
		if err != nil {
			log.Printf("rendering error: %v", err)
			c.SendString(fmt.Sprintf("rendering error: %v", err))
			return c.SendStatus(406)
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
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("T2P_PORT"))); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	_ = app.Shutdown()
}
