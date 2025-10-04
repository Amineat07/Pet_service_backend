package handler

import (
	"Pet_service_backend/db"
	"Pet_service_backend/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoute(app *fiber.App, con *pgxpool.Pool) {
	queries := db.New(con)

	auth := app.Group("/auth")
	auth.Post("/register", Register(queries))
	auth.Post("/login", Login(queries))

	service := app.Group("/service")
	service.Use(utils.JWTMiddleware([]byte(os.Getenv("JWT_SECRET")), queries))
	service.Post("/addservice", AddService(queries))
	service.Get("/services", GetServices(queries))
	service.Get("/service/:id", GetServicesByProvider(queries))
	service.Patch("services", UpdateServiceByProvider(queries))

	user := app.Group("/user")
	user.Use(utils.JWTMiddleware([]byte(os.Getenv("JWT_SECRET")), queries))
	user.Get("/users", GetUsers(queries))
	user.Patch("/user", UpdateUser(queries))
	user.Delete("/user", DeleteUser(queries))
	user.Get("/providers", GetProvider(queries))

	reservation := app.Group("/reservation")
	reservation.Use(utils.JWTMiddleware([]byte(os.Getenv("JWT_SECRET")), queries))
	reservation.Post("/", AddBooking(queries))

}
