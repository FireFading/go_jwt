package routes

import (
	"github.com/gofiber/fiber/v2"
	"go_jwt/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
}
