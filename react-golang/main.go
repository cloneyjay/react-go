package main

import (
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/react-golang/routes"
	"log"

	"github.com/gofiber/fiber"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)
	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
