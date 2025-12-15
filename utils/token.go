package utils

import (
	"Pet_service_backend/db"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(id uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"exp":     jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func JWTMiddleware(secret []byte, queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("jwt")

		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing token",
			})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		var userID int64
		switch id := claims["user_id"].(type) {
		case float64:
			userID = int64(id)
		case string:
			v, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "invalid user_id in token",
				})
			}
			userID = v
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}

		user, err := queries.GetUserById(c.UserContext(), userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user not found or deleted",
			})
		}

		c.Locals("user_id", user.ID)
		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}
		return c.Next()
	}
}
