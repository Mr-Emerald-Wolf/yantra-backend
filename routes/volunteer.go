package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func VolunteerRoutes(app *fiber.App) {
	app.Post("vol/create", controllers.CreateVolunteer)
	app.Get("vol/findall", controllers.GetVolunteers)
	app.Get("vol/find/:userId", controllers.FindVolByID)
	app.Post("vol/login", controllers.LoginVolunteer)
	app.Post("vol/refresh", controllers.RefreshToken)
	volGroup := app.Group("/vol", middleware.VerifyVol)
	volGroup.Get("/me", controllers.GetMeVolunteer)
	volGroup.Put("/update", controllers.UpdateVolunteer)
	volGroup.Delete("/delete", controllers.DeleteVolunteer)
}
