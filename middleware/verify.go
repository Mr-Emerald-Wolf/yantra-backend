package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/initializers"
	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/mr-emerald-wolf/yantra-backend/utils"
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

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Get the user
	findUser := models.User{}
	result := initializers.DB.First(&findUser, "email = ?", res.Email)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find user belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentUser", findUser.Email)
	return ctx.Next()
}

func VerifyAdmin(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "ADMIN" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an admin"})
	}
	// Get the user
	findUser := models.User{}
	result := initializers.DB.First(&findUser, "email = ?", res.Email)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find user belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentUser", findUser.Email)
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

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "NGO" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an NGO"})
	}
	// Get the user
	findNGO := models.Ngo{}
	result := initializers.DB.First(&findNGO, "email = ?", res.Email)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find NGO belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentNGO", findNGO.Email)
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

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "VOLUNTEER" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an NGO"})
	}
	// Get the user
	findVol := models.Volunteer{}
	result := initializers.DB.First(&findVol, "email = ?", res.Email)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error(), "message": "Could not find volunteer belonging to this token"})
		// log.Fatal(result.Error)
	}

	ctx.Set("currentVol", findVol.Email)
	return ctx.Next()
}
