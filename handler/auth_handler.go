package handler

import (
	"Pet_service_backend/db"
	requestresponse "Pet_service_backend/request_response"
	"Pet_service_backend/utils"
	"fmt"
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
		var req requestresponse.RegisterReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if err := utils.Validate(req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
		}

		fmt.Println("Incoming body:", string(c.Body()))

		if req.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please enter your password",
			})
		}

		if !passwordValidation(req.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please enter valid password",
			})
		}

		if !emailValidation(req.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please enter valid email",
			})
		}

		if *req.IsCustomer == true && *req.IsServiceProvider == true {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "You must choose one role only: Customer or Service Provider. Please make sure your selection is correct.",
			})
		}

		dbEmail, err := queries.CheckEmail(c.Context(), req.Email)
		if err == nil && req.Email == dbEmail {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}

		user, err := queries.CreateUser(c.Context(), db.CreateUserParams{
			Firstname:         req.Firstname,
			Lastname:          req.Lastname,
			Email:             req.Email,
			Iscustomer:        *req.IsCustomer,
			Isserviceprovider: *req.IsServiceProvider,
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
		var req requestresponse.LoginReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid login request",
			})
		}

		fmt.Println("body request", string(c.Body()))

		if err := utils.Validate(req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Validation error: %s", err))
		}

		user, err := queries.GetUserByEmail(c.Context(), req.Email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid email or password",
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
				"message": "Invalid email or password",
			})
		}

		token, err := utils.GenerateToken(uint(user.ID), role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to generate token",
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

		return c.JSON(requestresponse.LoginResponse{
			Firstname:  user.Firstname,
			Lastname:   user.Lastname,
			Email:      user.Email,
			Token:      token,
			IsCustomer: user.Iscustomer,
			IsProvider: user.Isserviceprovider,
			IsAdmin:    user.Isadmin,
		})
	}
}

func Logout(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HTTPOnly: true,
			Secure:   false,
			SameSite: fiber.CookieSameSiteStrictMode,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged out successfully",
		})
	}
}
