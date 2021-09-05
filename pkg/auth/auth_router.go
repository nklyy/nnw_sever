package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/user"
)

type Handler struct {
	authService IAuthService
	userService user.IUserService
	cfg         config.Configurations
	validate    *validator.Validate
}

func NewHandler(aService IAuthService,
	uService user.IUserService,
	cfg config.Configurations,
	v *validator.Validate) *Handler {
	return &Handler{
		authService: aService,
		userService: uService,
		cfg:         cfg,
		validate:    v,
	}
}

func (h *Handler) InitialRoute(route *echo.Echo) {
	v1 := route.Group("/v1")

	// Auth
	{
		// Registration and Verify Email
		v1.POST("/verifyRegistrationEmail", h.verifyRegistrationEmail)
		v1.POST("/verifyRegistrationEmailResend", h.verifyRegistrationEmailResend)
		v1.POST("/verifyRegistrationEmailCode", h.verifyRegistrationEmailCode)
		v1.POST("/verifyRegister2fa", h.verifyRegistration2FaCode)

		// Login
		v1.POST("/login", h.login)
		v1.POST("/verifyLogin2fa", h.verifyLogin2fa)

		v1.POST("/checkEmail", h.checkEmail)
		v1.POST("/checkJwt", h.checkJwt)
	}

	route.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
