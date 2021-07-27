package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"log"
	"net"
	"os"
)

func main() {
	// Initialise IP Database
	db := geoDB{}
	db.Setup()

	// Initialise fiber
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".gohtml"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("pages/index", fiber.Map{}, "templates/main")
	})

	app.Get("/github", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/maiacodes/geo")
	})

	app.Get("/:query", func(c *fiber.Ctx) error {
		query := c.Params("query")
		c.Set("Cache-Control", "public, max-age=604800")

		// Try IP
		ip := net.ParseIP(query)
		if ip == nil {
			n, err := net.LookupIP(query)
			if err != nil || len(n) == 0 {
				return c.Render("pages/error", fiber.Map{
					"IP": query,
				}, "templates/main")
			}
			ip = n[0]
		}

		// Query in DB
		info := db.Query(ip)
		return c.Render("pages/query", fiber.Map{
			"IP":   ip.String(),
			"Info": info,
		}, "templates/main")
	})

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
