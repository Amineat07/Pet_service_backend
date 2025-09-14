package main

import (
	"Pet_service_backend/database"
	"Pet_service_backend/handler"
	"Pet_service_backend/tutorial"

	"github.com/gofiber/fiber/v2"
)

func main() {

	con := database.ConnectDB()
	app := fiber.New()
	defer con.Close()
	queries := tutorial.New(con)

	app.Post("/register", handler.Register(queries))
	app.Post("/login", handler.Login(queries))

	app.Listen(":3000")
}
