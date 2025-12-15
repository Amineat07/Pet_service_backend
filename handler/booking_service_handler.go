package handler

import (
	"Pet_service_backend/db"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type BookingServiceRequest struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ServiceType string    `json:"service_type"`
	ProviderID  int64     `json:"provider_id"`
}

type BookingServiceResponse struct {
	ReservationID int64     `json:"reservation_id"`
	ProviderID    int64     `json:"provider_id"`
	ServiceType   string    `json:"service_type"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

func AddBookingService(queries *db.Queries) fiber.Handler {
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

		allowedServices := []string{
			"pet_sitting",
			"dog_walking",
			"pet_day_care",
			"pet_grooming",
			"pet_training",
			"pet_massage",
		}

		if req.ServiceType != "pet_sitting" &&
			req.ServiceType != "dog_walking" &&
			req.ServiceType != "pet_day_care" &&
			req.ServiceType != "pet_grooming" &&
			req.ServiceType != "pet_training" &&
			req.ServiceType != "pet_massage" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":            "invalid_service_type",
				"message":          "Unsupported service type. Please select one of the allowed options.",
				"allowed_services": allowedServices,
			})
		}

		customerID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}

		if !req.StartTime.Before(req.EndTime) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_time_range",
				"message": "start_time must be before end_time",
			})
		}

		if req.EndTime.Sub(req.StartTime) < 30*time.Minute {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "duration_too_short",
				"message": "Booking must be at least 30 minutes long",
			})
		}

		bookservice, err := queries.MakeReservation(c.Context(), db.MakeReservationParams{
			CustomerID:  customerID,
			ProviderID:  int64(req.ProviderID),
			ServiceType: req.ServiceType,
			StartTime:   pgtype.Timestamptz{Time: req.StartTime.UTC(), Valid: true},
			EndTime:     pgtype.Timestamptz{Time: req.EndTime.UTC(), Valid: true},
		})

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error":   "time_conflict",
					"message": "You already have a reservation during this time slot. Please choose a different time.",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "db_error",
				"message": "Failed to create reservation. Please try again later.",
			})
		}

		if bookservice.CustomerID == 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "time_conflict",
				"message": "You already have a reservation during this time slot.",
			})
		}

		response := BookingServiceResponse{
			ReservationID: bookservice.ID,
			ServiceType:   bookservice.ServiceType,
			ProviderID:    bookservice.ProviderID,
			StartTime:     bookservice.StartTime.Time.UTC(),
			EndTime:       bookservice.EndTime.Time.UTC(),
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}

}

type BookingServiceUpdateRequest struct {
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	ServiceType string     `json:"service_type"`
}

type BookingServiceUpdateResponse struct {
	ServiceType string    `json:"service_type"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

func UpdateBookingService(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "customer" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "acces denied: only customer can make booking!",
			})
		}

		var updateReq BookingServiceUpdateRequest
		if err := c.BodyParser(&updateReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid request body",
			})
		}

		if updateReq.StartTime == nil || updateReq.EndTime == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "start_time and end_time are required",
			})
		}

		idStr := c.Params("id")
		bookingID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid booking id",
			})
		}

		customerID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}

		booking, err := queries.GetReservation(c.Context(), bookingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "booking not found",
			})
		}

		if booking.CustomerID != customerID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "booking does not belong to user",
			})
		}

		providerServices, err := queries.GetServicesByProviderID(c.Context(), booking.ProviderID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to get provider services",
			})
		}

		var allowed bool
		switch updateReq.ServiceType {
		case "pet_sitting":
			allowed = providerServices.PetSitting
		case "dog_walking":
			allowed = providerServices.DogWalking
		case "pet_day_care":
			allowed = providerServices.PetDayCare
		case "pet_grooming":
			allowed = providerServices.PetGrooming
		case "pet_training":
			allowed = providerServices.PetTraining
		case "pet_massage":
			allowed = providerServices.PetMassage
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid service type",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("the provider does not offer the requested service: %s", updateReq.ServiceType),
			})
		}

		updated, err := queries.UpdateReservation(c.Context(), db.UpdateReservationParams{
			ID:        booking.ID,
			Column2:   updateReq.ServiceType,
			StartTime: pgtype.Timestamptz{Time: *updateReq.StartTime, Valid: true},
			EndTime:   pgtype.Timestamptz{Time: *updateReq.EndTime, Valid: true},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to update reservation",
			})
		}

		return c.Status(fiber.StatusOK).JSON(&updated)
	}
}

func DeleteReservation(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "customer" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "acces denied: only customer can make booking!",
			})
		}

		idStr := c.Params("id")
		bookingID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid booking id",
			})
		}

		if err := queries.DeleteReservation(c.Context(), bookingID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to delete reservation",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Reservation deleted.",
		})
	}
}
