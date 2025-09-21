package handler

import (
	"Pet_service_backend/db"
	"Pet_service_backend/utils"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UpdateUserReq struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}

func UpdateUser(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}
		var updateUserReq UpdateUserReq
		if err := c.BodyParser(&updateUserReq); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		current, err := queries.GetUserById(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "user not found"})
		}

		user := db.UpdateUserParams{
			ID:        userID,
			Firstname: utils.Pick(updateUserReq.FirstName, current.Firstname),
			Lastname:  utils.Pick(updateUserReq.LastName, current.Lastname),
			Email:     utils.Pick(updateUserReq.Email, current.Email),
			Password: func() string {
				if updateUserReq.Password != nil {
					return utils.GeneratePassword(*updateUserReq.Password)
				}
				return current.Password
			}(),
		}

		if err := queries.UpdateUser(c.Context(), user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(user)
	}
}

func DeleteUser(queries *db.Queries) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID, ok := c.Locals("user_id").(int64)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "invalid user_id in token",
            })
        }

        err := queries.DeleteUser(c.Context(), userID)
        if err != nil {
            if errors.Is(err, sql.ErrNoRows) {
                return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                    "message": "user not found",
                })
            }

            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "failed to delete user",
                "error":   err.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "message": "Your account has been deleted.",
        })
    }
}

