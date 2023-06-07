package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func NGORoutes(app *fiber.App) {
	app.Post("ngo/create", controllers.CreateNgo)
	app.Get("ngo/findall", controllers.AllNGOs)
	app.Get("ngo/find/:ngoId", controllers.FindNGObyId)
	app.Post("ngo/login", controllers.LoginNGO)
	app.Post("ngo/refresh", controllers.RefreshToken)
	ngoGroup := app.Group("/ngo", middleware.VerifyNGO)
	ngoGroup.Get("/me", controllers.GetNGO)
	ngoGroup.Put("/update", controllers.UpdateNGO)
	ngoGroup.Delete("/delete", controllers.DeleteNGO)
}
