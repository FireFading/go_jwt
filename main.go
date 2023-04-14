package main

import (
	"github.com/gofiber/fiber/v2"
	"go_jwt/database"
	"go_jwt/routes"
)

func main() {
	database.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":3000")
}
