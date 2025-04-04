package routes

import (
	"github.com/arezvani/wallet-go/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api")

	// Routes for GET method:
	route.Get("/transactions/:walletId", controllers.GetTransactions)
	route.Get("/balance/:walletId", controllers.GetBalance)
	route.Post("/transaction", controllers.PostTransaction)
}
