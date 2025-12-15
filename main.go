package main

import (
	"Pet_service_backend/database"
	"Pet_service_backend/handler"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	con := database.ConnectDB()
	app := fiber.New()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env")
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOWED_ORIGIN"),
		AllowCredentials: true,
	}))
	defer con.Close()

	handler.SetupRoute(app, con)

	app.Listen(":3000")
}
