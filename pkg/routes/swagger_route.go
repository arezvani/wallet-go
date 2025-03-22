package routes

import (
	"github.com/gofiber/fiber/v2"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/swagger-ui")

	// Routes
	route.Get("*", swagger.HandlerDefault)
	route.Post("*", swagger.HandlerDefault)
}
