package controllers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/cache"
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/mr-emerald-wolf/yantra-backend/utils"

	"gorm.io/gorm"
)

func CreateNgo(c *fiber.Ctx) error {
	var payload models.CreateNgoSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}
	hash, _ := utils.HashPassword(payload.Password)

	now := time.Now()
	newNgo := models.Ngo{
		Name:      payload.Name,
		Email:     payload.Email,
		Phone:     payload.Phone,
		Role:      payload.Role,
		Address:   payload.Address,
		Category:  payload.Category,
		Password:  hash,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newNgo)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Email already exist, please use another email"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "ngo": newNgo})
}

func AllNGOs(c *fiber.Ctx) error {

	var ngos []models.Ngo
	results := initializers.DB.Find(&ngos)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(ngos), "NGOs:": ngos})
}

func UpdateNGO(c *fiber.Ctx) error {
	email := c.GetRespHeader("currentNGO")

	var payload models.UpdateNgoSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var ngo models.Ngo
	result := initializers.DB.First(&ngo, "email = ?", email)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No NGO with that email exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.Phone != "" {
		updates["phone"] = payload.Phone
	}

	updates["updated_at"] = time.Now()

	fmt.Println(updates)

	initializers.DB.Model(&ngo).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"NGO": ngo}})
}

func FindNGObyId(c *fiber.Ctx) error {
	ngoId := c.Params("ngoId")
	// Get the NGO
	ngo := models.Ngo{}
	result := initializers.DB.First(&ngo, "id = ?", ngoId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "NGO": ngo})
}

func DeleteNGO(c *fiber.Ctx) error {
	email := c.GetRespHeader("currentNGO")
	// Get the ngo
	ngo := models.Ngo{}
	result := initializers.DB.First(&ngo, "email = ?", email)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "NGO Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.Ngo{}, "email = ?", email)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "NGO Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Ngo Deleted"})
}

func GetNGO(c *fiber.Ctx) error {
	email := c.GetRespHeader("currentNGO")
	ngo := models.Ngo{}

	result := initializers.DB.First(&ngo, "email = ?", fmt.Sprint(email))
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the NGO belonging to this token no logger exists"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "NGO": ngo})
}

func LoginNGO(c *fiber.Ctx) error {
	// Get request body and bind to payload
	var payload models.LoginNgoRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Validate Struct
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Get the NGO
	findNgo := models.Ngo{}
	result := initializers.DB.First(&findNgo, "email = ?", payload.Email)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "NGO Not Found"})
		// log.Fatal(result.Error)
	}
	// Compare password hashes
	match := utils.CheckPasswordHash(payload.Password, findNgo.Password)

	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Wrong password"})
	}

	// Create a new refreshToken
	duration, _ := time.ParseDuration("1h")
	sub := utils.TokenPayload{
		Email: findNgo.Email,
		Role:  findNgo.Role,
	}
	token, err := utils.GenerateToken(duration, sub, os.Getenv("REFRESH_JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	// Update refreshToken in NGO document
	errr := cache.SetValue(token, findNgo.Email, 0)
	if errr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "token": token})
}
