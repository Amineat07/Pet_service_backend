package main

import (
	"Pet_service_backend/database"
	"Pet_service_backend/handler"

	"github.com/gofiber/fiber/v2"
)

func main() {

	con := database.ConnectDB()
	app := fiber.New()
	defer con.Close()

	handler.SetupRoute(app, con)

	app.Listen(":3000")
}
