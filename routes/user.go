package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func UserRoutes(app *fiber.App) {
	app.Post("user/create", controllers.CreateUser)
	app.Get("user/findall", controllers.GetUsers)
	app.Get("user/find/:userId", controllers.FindUserById)
	app.Post("user/login", controllers.LoginUser)
	app.Post("user/refresh", controllers.RefreshToken)
	userGroup := app.Group("/user", middleware.VerifyUser)
	userGroup.Get("/me", controllers.GetMe)
	userGroup.Put("/update", controllers.UpdateUser)
	userGroup.Delete("/delete", controllers.DeleteUser)
}
