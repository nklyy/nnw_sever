package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/handler"
	"nnw_s/pkg/repository"
	"nnw_s/pkg/service"
	"os"
)

func Execute() {
	// Init config
	path := "."
	cfg, err := config.InitConfig(path, os.Getenv("APP_ENV"))
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	// Create App
	app := echo.New()

	// Connection to DB
	db, err := repository.MongoDbConnection(cfg)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	// Init App Middleware
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Set up validator
	validate := validator.New()

	// Init repository, service and handlers
	newRepository := repository.NewRepository(db, *cfg)
	newService := service.NewService(newRepository, *cfg)
	newHandler := handler.NewHandler(newService, *cfg, validate)

	newHandler.InitialRoute(app)

	// NotFound Urls
	echo.NotFoundHandler = func(c echo.Context) error {
		// Return HTTP 404 status and JSON response.
		return c.JSON(http.StatusNotFound, echo.Map{
			"error":    true,
			"endpoint": c.Request().URL.Path,
			"msg":      "Sorry, endpoint is not found",
		})
	}

	// Starting App
	err = app.Start(cfg.PORT)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}
