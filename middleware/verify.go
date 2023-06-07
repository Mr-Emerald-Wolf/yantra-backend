package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/mr-emerald-wolf/yantra-backend/utils"
	"github.com/spf13/viper"
)

func VerifyToken(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	_, err := utils.ValidateToken(token, viper.GetString("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return ctx.Next()
}

func VerifyUser(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, viper.GetString("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "USER" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Not a user"})
	}
	// Get the user
	findUser := models.User{}
	result := initializers.DB.First(&findUser, "id = ?", res.Id)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find user belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentUser", findUser.ID.String())
	return ctx.Next()
}

func VerifyNGO(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, viper.GetString("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "NGO" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an NGO"})
	}
	// Get the user
	findNGO := models.Ngo{}
	result := initializers.DB.First(&findNGO, "id = ?", res.Id)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find NGO belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentNGO", findNGO.ID.String())
	return ctx.Next()
}

func VerifyVol(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, viper.GetString("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "VOLUNTEER" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an Volunteer"})
	}
	// Get the user
	findVol := models.Volunteer{}
	result := initializers.DB.First(&findVol, "id = ?", res.Id)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find volunteer belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentVol", findVol.ID.String())
	return ctx.Next()
}
