package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func EventRoutes(app *fiber.App) {
	app.Get("event/findall", controllers.GetAllEvents)
	app.Get("event/find/:eventId", controllers.GetEvent)
	eventGroup := app.Group("/event", middleware.VerifyNGO)
	eventGroup.Post("/create", controllers.CreateEvent)
	eventGroup.Put("/update", controllers.UpdateBlog)
	eventGroup.Delete("/delete", controllers.DeleteUser)
}
