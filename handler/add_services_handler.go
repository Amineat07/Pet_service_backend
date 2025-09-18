package handler

import (
	"Pet_service_backend/db"

	"github.com/gofiber/fiber/v2"
)

type ServicesReq struct {
	ProviderID  bool `json:"provider_id"`
	PetSitting  bool `json:"pet_sitting"`
	DogWalking  bool `json:"dog_walking"`
	PetDayCare  bool `json:"pet_day_care"`
	PetGrooming bool `json:"pet_grooming"`
	PetTraining bool `json:"pet_training"`
	PetMassage  bool `json:"pet_massage"`
}

func AddService(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "provider" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "access denied: only providers or admins can add services",
			})
		}

		var req ServicesReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid request body",
			})
		}

		providerID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}

		service, err := queries.UpsertServices(c.Context(), db.UpsertServicesParams{
			ProviderID:  providerID,
			PetSitting:  req.PetSitting,
			DogWalking:  req.DogWalking,
			PetDayCare:  req.PetDayCare,
			PetGrooming: req.PetGrooming,
			PetTraining: req.PetTraining,
			PetMassage:  req.PetMassage,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(&service)
	}
}
