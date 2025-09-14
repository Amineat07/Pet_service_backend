package handler

import (
	reqres "Pet_service_backend/req-res"
	"Pet_service_backend/tutorial"
	"Pet_service_backend/utils"
	"regexp"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

func passwordValidation(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 12 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func emailValidation(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func Register(queries *tutorial.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req reqres.RegisterReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid request body",
			})
		}

		if req.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "please enter your password",
			})
		}

		if !passwordValidation(req.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "please enter valid password",
			})
		}

		if !emailValidation(req.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "please enter valid email",
			})
		}

		dbEmail, err := queries.CheckEmail(c.Context(), req.Email)
		if err == nil && req.Email == dbEmail {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "email already exists",
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

func Login(queries *tutorial.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req reqres.LoginReq
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

		return c.JSON(reqres.LoginResponse{
			Email: user.Email,
			Token: token,
		})
	}
}
