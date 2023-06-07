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

func CreateBlog(ctx *fiber.Ctx) error {
	var payload models.CreateBlogSchema

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	newBlog := models.Blog{
		Title:       payload.Title,
		SubTitle:    payload.SubTitle,
		Description: payload.Description,
		Author:      payload.Author,
		Date:        payload.Date,
	}

	result := initializers.DB.Create(&newBlog)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Blog already exist, please use another Blog name"})
	} else if result.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Blog created succesfully", "Blog": newBlog})
}

func GetAllBlogs(ctx *fiber.Ctx) error {
	var blogs []models.Blog
	results := initializers.DB.Find(&blogs)
	if results.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(blogs), "blogs": blogs})

}

func GetBlog(ctx *fiber.Ctx) error {
	BlogId := ctx.Params("blogId")
	// Get the Blog
	Blog := models.Blog{}
	result := initializers.DB.First(&Blog, "id = ?", BlogId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "Blog": Blog})
}

func UpdateBlog(ctx *fiber.Ctx) error {
	var payload models.UpdateBlogSchema
	BlogId := ctx.Params("blogId")

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var blog models.Blog
	result := initializers.DB.First(&blog, "id = ?", BlogId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No Blog with that email exists"})
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
	if *payload.Author != "" {
		updates["author"] = payload.Author
	}
	if *payload.Date != "" {
		updates["date"] = payload.Date
	}

	fmt.Println(updates)

	initializers.DB.Model(&blog).Updates(updates)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"blog": blog}})

}

func DeleteBlog(ctx *fiber.Ctx) error {

	blogId := ctx.Params("blogId")
	// Get the Blog
	blog := models.Blog{}
	result := initializers.DB.First(&blog, "id = ?", blogId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Blog Not Found", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	result = initializers.DB.Delete(&models.Blog{}, "id = ?", blogId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Blog Not Deleted", "err": result.Error.Error()})
		// log.Fatal(result.Error)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Blog Deleted"})

}
