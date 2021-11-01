package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"nnw_s/config"
	"nnw_s/internal/auth"
	"nnw_s/internal/auth/jwt"
	"nnw_s/internal/auth/twofa"
	"nnw_s/internal/auth/verification"
	"nnw_s/internal/user"
	"nnw_s/internal/user/credentials"
	"nnw_s/internal/wallet"
	"nnw_s/pkg/mongodb"
	"nnw_s/pkg/notificator"
	"nnw_s/pkg/smtp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init config
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	logger := logrus.New()

	// Create App
	router := echo.New()

	// Connection to DB
	db, err := mongodb.NewConn(cfg)
	if err != nil {
		logger.Fatalf("failed to connect to mongodb: %v", err)
	}

	// Init App Middleware
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.CorsOrigin.DevOrigin, cfg.CorsOrigin.ProdOrigin},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Init dependencies
	userRepo := user.NewRepository(db, logger)
	credentialsSvc, err := credentials.NewService(logger, cfg.Shift, cfg.PasswordSalt)
	if err != nil {
		logger.Fatalf("failed to create credentials service: %v", err)
	}

	userSvc, err := user.NewService(userRepo, credentialsSvc, logger)
	if err != nil {
		logger.Fatalf("failed to create user service: %v", err)
	}

	smtpClient := smtp.NewClient(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUserApiKey, cfg.SmtpPasswordKey)
	notificatorSvc, err := notificator.NewService(logger, smtpClient)
	if err != nil {
		logger.Fatalf("failed to create user service: %v", err)
	}

	verificationRepo, err := verification.NewRepository(db, logger)
	if err != nil {
		logger.Fatalf("failed to create verification repo: %v", err)
	}

	verificationSvc, err := verification.NewService(verificationRepo, logger)
	if err != nil {
		logger.Fatalf("failed to create verification service: %v", err)
	}

	twoFaSvc, err := twofa.NewService(cfg.TwoFAIssuer)
	if err != nil {
		logger.Fatalf("failed to create TwoFA service: %v", err)
	}

	jwtRepo, err := jwt.NewRepository(db)
	if err != nil {
		logger.Fatalf("failed to create JWT repo: %v", err)
	}

	jwtSvc, err := jwt.NewService(jwtRepo, cfg.JwtSecretKey)
	if err != nil {
		logger.Fatalf("failed to create JWT service: %v", err)
	}

	authDeps := auth.ServiceDeps{
		UserService:         userSvc,
		NotificatorService:  notificatorSvc,
		VerificationService: verificationSvc,
		TwoFAService:        twoFaSvc,
		JWTService:          jwtSvc,
		CredentialsService:  credentialsSvc,
	}

	registrationSvc, err := auth.NewRegistrationService(logger, cfg.EmailFrom, &authDeps)
	if err != nil {
		logger.Fatalf("failed to create registration service: %v", err)
	}

	loginSvc, err := auth.NewLoginService(logger, &authDeps)
	if err != nil {
		logger.Fatalf("failed to connect login service: %v", err)
	}

	resetPasswordSvc, err := auth.NewResetPasswordService(logger, cfg.EmailFrom, &authDeps)
	if err != nil {
		logger.Fatalf("failed to connect reset password service: %v", err)
	}

	walletDeps := wallet.ServiceDeps{
		UserService:        userSvc,
		TwoFAService:       twoFaSvc,
		JWTService:         jwtSvc,
		CredentialsService: credentialsSvc,
	}

	walletSvc, err := wallet.NewWalletService(logger, &walletDeps)
	if err != nil {
		logger.Fatalf("failed to connect wallet wallet service: %v", err)
	}

	// Handlers
	// User
	userHandler := user.NewHandler(userSvc, jwtSvc, cfg.Shift)
	userHandler.SetupRoutes(router)

	// Auth
	authHandler := auth.NewHandler(registrationSvc, loginSvc, resetPasswordSvc, jwtSvc, cfg.Shift)
	authHandler.SetupRoutes(router)

	// Wallet
	walletHandler := wallet.NewHandler(walletSvc, jwtSvc, cfg.Shift)
	walletHandler.SetupRoutes(router)

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
	logger.Info("starting NNW server at :%s...", cfg.PORT)
	if err = router.Start(":" + cfg.PORT); err != nil {
		logger.Errorf("failed to start HTTP server: %v", err)
		return
	}
}
