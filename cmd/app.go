package cmd

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"nnw_s/config"
	"nnw_s/pkg/repository"
)

func Execute() {
	// Init config
	cfg, err := initConfig()
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	// Create App
	app := fiber.New()

	// Connection to DB
	_, err = repository.MongoDbConnection(cfg)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	//// Init repository, service and handlers
	//newRepository := repository.NewRepository(db)
	//newService := service.NewService(newRepository)
	//newHandler := handler.NewHandler(newService)
	//
	//newHandler.InitialRoute(app)

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

func initConfig() (*config.Configurations, error) {
	viper.AddConfigPath("config")

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configuration config.Configurations
	err = viper.Unmarshal(&configuration)
	if err != nil {
		//fmt.Printf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &configuration, nil
}
