package routes

import (
	"docs_api_golang/controllers"

	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {
	// app.Get("/users", handlers.getUser)
	// app.Get("/users/:id?", handlers.getUser)
	// app.Put("/users/:id", handlers.updateUser)
	app.Post("/users", controllers.CreateUser)
	// app.Delete("/users/:id", handlers.deleteUser)

	// app.Get("/documents", handlers.getDocument)
	// app.Get("/documents/:id?", handlers.getDocument)
	// app.Put("/documents/:id", handlers.updateDocument)
	app.Post("/documents", controllers.CreateDocument)
	// app.Delete("/documents/:id", handlers.deleteDocument)
}
