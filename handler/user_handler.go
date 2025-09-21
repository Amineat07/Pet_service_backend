package handler

import (
	"Pet_service_backend/db"
	requestresponse "Pet_service_backend/request_response"
	"Pet_service_backend/utils"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dbusers, err := queries.GetUsers(c.Context())
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		users := make([]requestresponse.UsersResponse, len(dbusers))
		for i, u := range dbusers {
			users[i] = requestresponse.UsersResponse{
				FirstName:  u.Firstname,
				LastName:   u.Lastname,
				Email:      u.Email,
				IsCustomer: u.Iscustomer,
				IsProvider: u.Isserviceprovider,
				Created_At: u.CreatedAt.Time,
			}
		}
		return c.JSON(users)
	}
}

func GetProvider(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {

		role, _ := c.Locals("role").(string)
		if role != "customer" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "access denied: only customers or admins can see services",
			})
		}

		providersdb, err := queries.GetProviders(c.Context())
		if err != nil {
			fmt.Printf("GetProviders error: %v\n", err) // Error Debuging
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		providers := make([]requestresponse.ProviderResponse, len(providersdb))
		for i, u := range providersdb {
			providers[i] = requestresponse.ProviderResponse{
				FirstName:  u.Firstname,
				LastName:   u.Lastname,
				Email:      u.Email,
				IsProvider: u.Isserviceprovider,
				Created_At: u.CreatedAt.Time,
			}
		}
		return c.JSON(providers)
	}
}

func UpdateUser(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}
		var updateUserReq requestresponse.UpdateUserReq
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

		if err := queries.DeleteServices(c.UserContext(), userID); err != nil {
			fmt.Printf("DeleteServices error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to delete user services",
				"error":   err.Error(),
			})
		}

		if err := queries.DeleteUser(c.UserContext(), userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "user not found",
				})
			}
			fmt.Printf("DeleteUser error: %v\n", err)
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
