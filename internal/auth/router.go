package auth

import (
	"net/http"
	"nnw_s/config"
	"nnw_s/internal/user"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService Service
	userService user.Service
	cfg         config.Config
	validate    *validator.Validate
}

func NewHandler(authService Service,
	userService user.Service,
	cfg config.Config,
	v *validator.Validate) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
		cfg:         cfg,
		validate:    v,
	}
}

func (h *Handler) InitialRoute(route *echo.Echo) {
	v1 := route.Group("/v1")
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

	route.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
