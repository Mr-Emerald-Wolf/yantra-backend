package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/cache"
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/mr-emerald-wolf/yantra-backend/utils"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

func CreateVolunteer(c *fiber.Ctx) error {
	var payload models.CreateVolunteerSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}
	hash, _ := utils.HashPassword(payload.Password)

	now := time.Now()
	newUser := models.Volunteer{
		Name:      payload.Name,
		Email:     payload.Email,
		Phone:     payload.Phone,
		Role:      "VOLUNTEER",
		Gender:    payload.Gender,
		Address:   payload.Address,
		Password:  hash,
		Category:  payload.Category,
		NGO:       payload.NGO,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Email already exist, please use another email"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "user": newUser})
}

func GetVolunteers(c *fiber.Ctx) error {

	var volunteers []models.Volunteer
	results := initializers.DB.Find(&volunteers)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(volunteers), "volunteers": volunteers})
}

func UpdateVolunteer(c *fiber.Ctx) error {
	id := c.GetRespHeader("currentVol")

	var payload *models.UpdateVolunteerSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var vol models.Volunteer
	result := initializers.DB.First(&vol, "id = ?", id)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No volunteer with that email exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.NGO.String() != "" {
		updates["NGO"] = payload.NGO
	}
	if payload.Phone != "" {
		updates["phone"] = payload.Phone
	}

	updates["updated_at"] = time.Now()

	fmt.Println(updates)

	initializers.DB.Model(&vol).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"volunteer": vol}})
}

func FindVolByID(c *fiber.Ctx) error {
	volId := c.Params("volId")
	// Get the user
	vol := models.Volunteer{}
	result := initializers.DB.First(&vol, "id = ?", volId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "volunteer": volId})
}

func DeleteVolunteer(c *fiber.Ctx) error {
	id := c.GetRespHeader("currentVol")
	// Get the volunteer
	vol := models.Volunteer{}
	result := initializers.DB.First(&vol, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.Volunteer{}, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User Deleted"})
}

func GetMeVolunteer(c *fiber.Ctx) error {
	id := c.GetRespHeader("currentVol")
	vol := models.Volunteer{}

	result := initializers.DB.First(&vol, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "volunteer": vol})
}

func LoginVolunteer(c *fiber.Ctx) error {
	// Get request body and bind to payload
	var payload *models.LoginVolunteerRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Validate Struct
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Get the user
	findVol := models.Volunteer{}
	result := initializers.DB.First(&findVol, "email = ?", payload.Email)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "Volunteer Not Found"})
		// log.Fatal(result.Error)
	}
	// Compare password hashes
	match := utils.CheckPasswordHash(payload.Password, findVol.Password)

	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Wrong password"})
	}

	// Create a new refreshToken
	duration, _ := time.ParseDuration("1h")
	sub := utils.TokenPayload{
		Id:   findVol.ID,
		Role: findVol.Role,
	}

	refreshSecret := viper.GetString("REFRESH_JWT_SECRET")
	token, err := utils.GenerateToken(duration, sub, refreshSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	// Update refreshToken in user document
	errr := cache.SetValue(token, findVol.Email, 0)
	if errr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "token": token})
}
