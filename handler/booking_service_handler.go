package handler

import (
	"Pet_service_backend/db"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type BookingServiceRequest struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ServiceType string    `json:"service_type"`
	ProviderID  uint      `json:"provider_id"`
}

func AddBooking(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "customer" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "acces denied: only customer can make booking!",
			})
		}

		var req BookingServiceRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid body request",
			})
		}

		customerID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}

		bookservice, err := queries.MakeReservation(c.Context(), db.MakeReservationParams{
			CustomerID:  customerID,
			ProviderID:  int64(req.ProviderID),
			ServiceType: req.ServiceType,
			StartTime:   pgtype.Timestamptz{Time: req.StartTime, Valid: true},
			EndTime:     pgtype.Timestamptz{Time: req.EndTime, Valid: true},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(&bookservice)
	}
}
