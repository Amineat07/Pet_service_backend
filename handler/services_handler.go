package handler

import (
	"Pet_service_backend/db"
	requestresponse "Pet_service_backend/request_response"
	"Pet_service_backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddService(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "provider" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "access denied: only providers or admins can add services",
			})
		}

		var req requestresponse.ServicesReq
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

func GetServices(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		services, err := queries.GetServices(c.Context())
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(services)
	}
}

func GetServicesByProvider(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
		proviedrIDParam := c.Params("id")
		serviceID, err := strconv.ParseInt(proviedrIDParam, 10, 64)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		service, err := queries.GetServiceByProviderID(c.Context(), serviceID)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(service)
	}
}

func UpdateServiceByProvider(queries *db.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {

		role, _ := c.Locals("role").(string)
		if role != "provider" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "access denied: only providers or admins can add services",
			})
		}

		providerID, ok := c.Locals("user_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid user_id in token",
			})
		}

		var updateServiceReq requestresponse.UpdateServicesReq
		if err := c.BodyParser(&updateServiceReq); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		current, err := queries.GetServiceByProviderID(c.Context(), providerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "user not found"})
		}

		services := db.UpdateServicesParams{
			PetSitting:  utils.Pick(updateServiceReq.PetSitting, current.PetSitting),
			DogWalking:  utils.Pick(updateServiceReq.DogWalking, current.DogWalking),
			PetDayCare:  utils.Pick(updateServiceReq.PetDayCare, current.PetDayCare),
			PetGrooming: utils.Pick(updateServiceReq.PetGrooming, current.PetGrooming),
			PetTraining: utils.Pick(updateServiceReq.PetTraining, current.PetTraining),
			PetMassage:  utils.Pick(updateServiceReq.PetMassage, current.PetMassage),
			ProviderID:  providerID,
		}

		if err := queries.UpdateServices(c.Context(), services); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(services)
	}
}
