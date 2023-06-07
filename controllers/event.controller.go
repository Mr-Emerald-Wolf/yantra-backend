package controllers

import (
	"fmt"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/mr-emerald-wolf/yantra-backend/utils"
	"gorm.io/gorm"
)

func CreateEvent(ctx *fiber.Ctx) error {
	var payload *models.Event

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	newEvent := models.Event{
		Title:       payload.Title,
		SubTitle:    payload.SubTitle,
		Description: payload.Description,
		Location:    payload.Location,
		Date:        payload.Date,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		NGO:         payload.NGO,
	}

	result := initializers.DB.Create(&newEvent)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Event already exist, please use another event name"})
	} else if result.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Event created succesfully", "Event": newEvent})
}

func GetAllEvents(ctx *fiber.Ctx) error {
	var events []models.Event
	results := initializers.DB.Find(&events)
	if results.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(events), "events": events})

}

func GetEvent(ctx *fiber.Ctx) error {
	eventId := ctx.Params("eventId")
	// Get the Event
	event := models.Event{}
	result := initializers.DB.First(&event, "id = ?", eventId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "event": event})
}

func UpdateEvent(ctx *fiber.Ctx) error {
	var payload models.UpdateEventSchema
	eventId := ctx.Params("eventId")

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var event models.Event
	result := initializers.DB.First(&event, "id = ?", eventId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No event with that email exists"})
		}
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if *payload.Title != "" {
		updates["title"] = payload.Title
	}
	if *payload.SubTitle != "" {
		updates["subtitle"] = payload.SubTitle
	}
	if *payload.Description != "" {
		updates["description"] = payload.Description
	}
	if *payload.Location != "" {
		updates["location"] = payload.Location
	}
	if *payload.Date != "" {
		updates["date"] = payload.Date
	}
	if *payload.StartTime != "" {
		updates["starttime"] = payload.StartTime
	}
	if *payload.EndTime != "" {
		updates["endtime"] = payload.EndTime
	}

	fmt.Println(updates)

	initializers.DB.Model(&event).Updates(updates)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"event": event}})

}

func DeleteEvent(ctx *fiber.Ctx) error {

	eventId := ctx.Params("eventId")
	// Get the event
	event := models.Event{}
	result := initializers.DB.First(&event, "id = ?", eventId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Event Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.Event{}, "id = ?", eventId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Event Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Event Deleted"})

}
