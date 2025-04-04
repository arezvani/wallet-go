package main

import (
	"github.com/arezvani/wallet-go/internal/configs"
	"github.com/arezvani/wallet-go/internal/database"
	"github.com/arezvani/wallet-go/internal/routes"
	"github.com/arezvani/wallet-go/pkg/middleware"
	"github.com/arezvani/wallet-go/pkg/utils"

	_ "github.com/arezvani/wallet-go/docs" // load API Docs files (Swagger)

	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically

	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v4/stdlib"
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
	// Create database connection
	if err := database.OpenDBConnection(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app)

	// Routes.
	routes.SwaggerRoute(app)
	routes.PublicRoutes(app)

	utils.StartServer(app)
}
