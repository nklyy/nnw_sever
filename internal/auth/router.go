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
	v1.POST("/preRegistration", h.preRegistration)
	v1.POST("/sendVerifyRegistrationEmail", h.sendVerifyRegistrationEmail)
	v1.POST("/resendVerifyRegistrationEmail", h.resendVerifyRegistrationEmail)
	v1.POST("/checkRegistrationEmailCode", h.checkRegistrationEmailCode)
	v1.POST("/generateRegistration2FA", h.generate2FAQrCode)
	v1.POST("/checkRegistration2FA", h.checkRegistration2FaCode)
	v1.POST("/finishRegistration", h.finishRegistration)

	// Login
	v1.POST("/login", h.login)
	v1.POST("/checkLogin2fa", h.checkLogin2faCode)

	v1.POST("/checkEmail", h.checkEmail)
	v1.POST("/checkJwt", h.checkJwt)

	route.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
