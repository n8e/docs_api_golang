package routes

import (
	"docs_api_golang/controllers"

	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {
	app.Post("/auth", controllers.Auth)

	app.Get("/users", controllers.GetUsers)
	app.Get("/users/:id", controllers.GetUser)
	app.Put("/users/:id", controllers.UpdateUser)
	app.Post("/users", controllers.CreateUser)
	app.Delete("/users/:id", controllers.DeleteUser)

	app.Get("/documents", controllers.GetDocuments)
	app.Get("/documents/:id?", controllers.GetDocument)
	app.Put("/documents/:id", controllers.UpdateDocument)
	app.Post("/documents", controllers.CreateDocument)
	app.Delete("/documents/:id", controllers.DeleteDocument)
}
