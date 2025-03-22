package main

import (
	"github.com/arezvani/wallet-go/pkg/configs"
	"github.com/arezvani/wallet-go/pkg/middleware"
	"github.com/arezvani/wallet-go/pkg/routes"
	"github.com/arezvani/wallet-go/utils"

	_ "github.com/arezvani/wallet-go/docs" // load API Docs files (Swagger)

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	_ "github.com/lib/pq"
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fcairib76@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app) // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app) // Register a public routes for app.

	utils.StartServer(app)
}
