package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/auth/handler"
	repository1 "nnw_s/pkg/auth/repository"
	service1 "nnw_s/pkg/auth/service"
	"nnw_s/pkg/common"
	repository2 "nnw_s/pkg/user/repository"
	service2 "nnw_s/pkg/user/service"
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
	newAuthRepository := repository1.NewRepository(db, *cfg)
	newUserRepository := repository2.NewRepository(db, *cfg)
	newAuthService := service1.NewService(newAuthRepository, newUserRepository, *cfg)
	newUserService := service2.NewUserService(newUserRepository, *cfg)
	newAuthHandler := handler.NewHandler(newAuthService, newUserService, *cfg, validate)

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
