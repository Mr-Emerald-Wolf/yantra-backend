package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func VolunteerRoutes(app *fiber.App) {
	app.Post("vol/create", controllers.CreateUser)
	app.Get("vol/findall", controllers.GetUsers)
	app.Get("vol/find/:userId", controllers.FindUserById)
	app.Post("vol/login", controllers.LoginUser)
	app.Post("vol/refresh", controllers.RefreshToken)
	volGroup := app.Group("/vol", middleware.VerifyVol)
	volGroup.Get("/me", controllers.GetMe)
	volGroup.Put("/update", controllers.UpdateUser)
	volGroup.Delete("/delete", controllers.DeleteUser)
}
