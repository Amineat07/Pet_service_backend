package handler

import (
	"Pet_service_backend/tutorial"
	"Pet_service_backend/utils"

	"github.com/gofiber/fiber/v2"
)

type RegisterReq struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func Register(queries *tutorial.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid request body",
			})
		}

		user, err := queries.CreateUser(c.Context(), tutorial.CreateUserParams{
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Email:     req.Email,
			Password:  utils.GeneratePassword(req.Password),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(&user)

	}
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(queries *tutorial.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid login request",
			})
		}

		user, err := queries.GetUserByEmail(c.Context(), req.Email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid email or password",
			})
		}

		if !utils.ComparePassword(user.Password, req.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid email or password",
			})
		}

		token, err := utils.GenerateToken(uint(user.ID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate token",
			})
		}

		return c.JSON(fiber.Map{
			"email": user.Email,
			"token": token,
		})
	}
}
