package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-emerald-wolf/yantra-backend/controllers"
	"github.com/mr-emerald-wolf/yantra-backend/middleware"
)

func RequestRoutes(app *fiber.App) {
	app.Get("request/findall", controllers.AllNGOs)
	requestGroup := app.Group("/request", middleware.VerifyUser)
	requestGroup.Post("/create", controllers.CreateRequest)
	requestGroup.Get("/find/:requestId", controllers.FindNGObyId)
	requestGroup.Get("/fulfilled", controllers.GetFulfilledRequest)
	requestGroup.Get("/unfulfilled", controllers.GetUnFulfilledRequest)
	requestGroup.Put("/fulfill/:requestId", controllers.FulfillRequest)
	ngoGroup := app.Group("/ngo", middleware.VerifyNGO)
	ngoGroup.Get("/request", controllers.GetNgoRequests)
	volGroup := app.Group("/volunteer", middleware.VerifyVol)
	volGroup.Get("/request", controllers.GetVolRequests)

	// requestGroup.Delete("/delete", controllers.DeleteNGO)
}
