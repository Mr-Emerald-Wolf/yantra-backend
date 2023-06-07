package controllers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/mr-emerald-wolf/yantra-backend/utils"
)

func CreateRequest(c *fiber.Ctx) error {
	var payload *models.CreateRequestSchema
	id := c.GetRespHeader("currentUser")
	uuid, errr := uuid.Parse(id)
	if errr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": errr.Error()})
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	now := time.Now()
	newUser := models.Request{
		UserID:      uuid,
		Title:       payload.Title,
		Description: payload.Description,
		IsFulfilled: false,
		VolunteerID: payload.VolunteerID,
		NGO:         payload.NGO,
		CreatedAt:   now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Request already exist, please use another id"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "request": newUser})
}

func GetFulfilledRequest(c *fiber.Ctx) error {
	var requests []models.Request
	userId := c.GetRespHeader("currentUser")

	results := initializers.DB.Find(&requests, "userId = ? AND is_fulfilled = ?", userId, true)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(requests), "requests": requests})
}
