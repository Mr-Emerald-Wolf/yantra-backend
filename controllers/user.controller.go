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

func CreateUser(c *fiber.Ctx) error {
	var payload *models.CreateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}
	hash, _ := utils.HashPassword(payload.Password)

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Phone:     payload.Phone,
		Role:      payload.Role,
		Gender:    payload.Gender,
		Address:   payload.Address,
		Password:  hash,
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

func GetUsers(c *fiber.Ctx) error {

	var users []models.User
	results := initializers.DB.Find(&users)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(users), "users": users})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.GetRespHeader("currentUser")

	var payload *models.UpdateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := initializers.DB.First(&user, "id = ?", id)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that email exists"})
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

	initializers.DB.Model(&user).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")
	// Get the user
	user := models.User{}
	result := initializers.DB.First(&user, "id = ?", userId)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.GetRespHeader("currentUser")
	// Get the user
	user := models.User{}
	result := initializers.DB.First(&user, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.User{}, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User Deleted"})
}

func GetMe(c *fiber.Ctx) error {
	id := c.GetRespHeader("currentUser")
	user := models.User{}

	result := initializers.DB.First(&user, "id = ?", fmt.Sprint(id))
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": user})
}

func LoginUser(c *fiber.Ctx) error {
	// Get request body and bind to payload
	var payload *models.LoginUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Validate Struct
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Get the user
	findUser := models.User{}
	result := initializers.DB.First(&findUser, "email = ?", payload.Email)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "User Not Found"})
		// log.Fatal(result.Error)
	}
	// Compare password hashes
	match := utils.CheckPasswordHash(payload.Password, findUser.Password)

	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Wrong password"})
	}

	// Create a new refreshToken
	duration, _ := time.ParseDuration("1h")
	sub := utils.TokenPayload{
		Id: findUser.ID,
		Role:  findUser.Role,
	}

	refreshSecret := viper.GetString("REFRESH_JWT_SECRET")
	token, err := utils.GenerateToken(duration, sub, refreshSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	// Update refreshToken in user document
	errr := cache.SetValue(token, findUser.Email, 0)
	if errr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "token": token})
}

func RefreshToken(ctx *fiber.Ctx) error {

	payload := utils.TokenRequest{}
	// Get refreshToken from request
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Find refresh token in redis
	_, err := cache.GetValue(payload.RefreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "refreshToken not found"})
	}

	// Validate Refresh Token
	refreshSecret := viper.GetString("REFRESH_JWT_SECRET")
	sub, err := utils.ValidateToken(payload.RefreshToken, refreshSecret)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "err": err.Error(), "token": payload.RefreshToken})
	}

	// Create new accessToken
	duration, _ := time.ParseDuration("1h")
	secret := viper.GetString("JWT_SECRET")
	accessToken, err := utils.GenerateToken(duration, sub, secret)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "accessToken": accessToken})
}
