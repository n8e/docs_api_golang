package main

import (
	"docs_api_golang/configs"
	"docs_api_golang/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// run database
	configs.ConnectDB()

	// routes
	routes.AllRoutes(app)

	app.Listen(":8000")
}
