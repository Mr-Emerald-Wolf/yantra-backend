package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	UserRoutes(app)
	NGORoutes(app)
	VolunteerRoutes(app)
	EventRoutes(app)
	BlogRoutes(app)
}
