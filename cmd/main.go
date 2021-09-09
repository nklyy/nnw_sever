package cmd

import (
	"fmt"
	"log"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/mongodb"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init config
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create App
	app := echo.New()

	// Connection to DB
	db, err := mongodb.NewConn(cfg)
	if err != nil {
		fmt.Printf("ERROR: %s \n", err)
		return
	}

	// Init App Middleware
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Set-up SMTP
	smtpPort, err := strconv.Atoi(cfg.SmtpPort)
	if err != nil {
		fmt.Println("ERROR: Incorrect SMTP PORT!")
		return
	}
	emailClient := config.NewSMTPClient(cfg.SmtpHost, smtpPort, cfg.SmtpUserApiKey, cfg.SmtpPasswordKey)

	// Set up validator
	validate := validator.New()

	// Init repository, service and handlers
	newAuthRepository := auth.NewAuthRepository(db, *cfg)
	newUserRepository := user.NewUserRepository(db, *cfg)
	newAuthService := auth.NewAuthService(*newAuthRepository, *newUserRepository, *cfg, *emailClient)
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
