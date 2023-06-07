package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func BlogRoutes(app *fiber.App) {
	app.Get("blog/findall", controllers.GetAllBlogs)
	app.Get("blog/find/:blogId", controllers.GetBlog)
	blogGroup := app.Group("/blog", middleware.VerifyToken)
	blogGroup.Post("/create", controllers.CreateBlog)
	blogGroup.Put("/update", controllers.UpdateBlog)
	blogGroup.Delete("/delete", controllers.DeleteBlog)
}
