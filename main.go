package main

import (
	"docs_api_golang/configs"
	"docs_api_golang/controllers"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()

	// run database
	configs.ConnectDB()

	app.Post("/auth", controllers.Auth)
	app.Post("/users", controllers.CreateUser)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(configs.EnvVariable("JWT_SECRET_KEY")),
	}))

	// routes
	app.Get("/users", controllers.GetUsers)
	app.Get("/users/:id", controllers.GetUser)
	app.Put("/users/:id", controllers.UpdateUser)
	app.Delete("/users/:id", controllers.DeleteUser)

	app.Get("/documents", controllers.GetDocuments)
	app.Get("/documents/:id?", controllers.GetDocument)
	app.Put("/documents/:id", controllers.UpdateDocument)
	app.Post("/documents", controllers.CreateDocument)
	app.Delete("/documents/:id", controllers.DeleteDocument)

	app.Listen(":8000")
}
