package handler

import (
	"Pet_service_backend/db"
	reqres "Pet_service_backend/req-res"
	"Pet_service_backend/utils"
	"regexp"
	"time"
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

func Register(queries *db.Queries) fiber.Handler {
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

		if req.IsCustomer == true && req.IsServiceProvider == true {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "You must choose one role only: Customer or Service Provider. Please make sure your selection is correct.",
			})
		}

		dbEmail, err := queries.CheckEmail(c.Context(), req.Email)
		if err == nil && req.Email == dbEmail {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "email already exists",
			})
		}

		user, err := queries.CreateUser(c.Context(), db.CreateUserParams{
			Firstname:         req.Firstname,
			Lastname:          req.Lastname,
			Email:             req.Email,
			Iscustomer:        req.IsCustomer,
			Isserviceprovider: req.IsServiceProvider,
			Password:          utils.GeneratePassword(req.Password),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(&user)

	}
}

func Login(queries *db.Queries) fiber.Handler {
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

		var role string
		if user.Isadmin {
			role = "admin"
		} else if user.Iscustomer {
			role = "customer"
		} else if user.Isserviceprovider {
			role = "provider"
		}

		if !utils.ComparePassword(user.Password, req.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid email or password",
			})
		}

		token, err := utils.GenerateToken(uint(user.ID), role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to generate token",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
		})

		return c.JSON(reqres.LoginResponse{
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Token:     token,
		})
	}
}
