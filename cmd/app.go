package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/auth"
	"nnw_s/pkg/common"
	"nnw_s/pkg/user"
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
	db, err := common.MongoDbConnection(cfg)
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
	newAuthRepository := auth.NewAuthRepository(db, *cfg)
	newUserRepository := user.NewUserRepository(db, *cfg)
	newAuthService := auth.NewAuthService(*newAuthRepository, *newUserRepository, *cfg)
	newUserService := user.NewUserService(newUserRepository, *cfg)
	newAuthHandler := auth.NewHandler(newAuthService, newUserService, *cfg, validate)

	newAuthHandler.InitialRoute(app)

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
	err = app.Start(":" + cfg.PORT)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}
}
