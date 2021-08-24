package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/service"
)

type Handler struct {
	services *service.Service
	cfg      config.Configurations
}

func NewHandler(services *service.Service, cfg config.Configurations) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) InitialRoute(route *echo.Echo) {
	v1 := route.Group("/v1")

	// Auth
	{
		// Registration
		v1.POST("/registration", h.registration)
		v1.POST("/verifyRegister2fa", h.verifyRegistration2FaCode)

		// Login
		v1.POST("/login", h.login)
		v1.POST("/verifyLogin2fa", h.verifyLogin2fa)

		v1.POST("/checkLogin", h.checkUserName)
		v1.POST("/checkJwt", h.checkJwt)
	}

	route.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
