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

	// Init App Middleware
	app.Use(
		// Add CORS to each route.
		cors.New(),
	)

	//// Init repository, service and handlers
	newRepository := repository.NewRepository(db, *cfg)
	newService := service.NewService(newRepository, *cfg)
	newHandler := handler.NewHandler(newService, *cfg)

	newHandler.InitialRoute(app)

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
