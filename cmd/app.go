package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"nnw_s/config"
	"nnw_s/pkg/handler"
	"nnw_s/pkg/repository"
	"nnw_s/pkg/service"
)

func Execute() {
	// Init config
	path := "config"
	cfg, err := config.InitConfig(path)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	// Create App
	app := fiber.New()

	// Connection to DB
	db, err := repository.MongoDbConnection(cfg)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	//// Init repository, service and handlers
	newRepository := repository.NewRepository(db)
	newService := service.NewService(newRepository)
	newHandler := handler.NewHandler(newService)

	newHandler.InitialRoute(app)

	// Init App Middleware
	app.Use(
		// Add CORS to each route.
		cors.New(),
	)

	// NotFound Urls
	app.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "Sorry, endpoint " + "'" + c.OriginalURL() + "'" + " is not found",
			})
		},
	)

	// Starting App
	err = app.Listen(cfg.PORT)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}
