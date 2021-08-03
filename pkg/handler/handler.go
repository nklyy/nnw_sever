package handler

import (
	"github.com/gofiber/fiber/v2"
	"nnw_s/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitialRoute(route fiber.Router) {
	v1 := route.Group("/v1")

	// Auth
	{
		// Registration
		v1.Post("/registration", h.registration)
		v1.Post("/verifyRegister2fa", h.verifyRegistration2FaCode)

		// Login
		v1.Post("/login", h.login)
		v1.Post("/verifyLogin2fa", h.verifyLogin2fa)

		v1.Post("/checkLogin", h.checkLogin)
	}

	route.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("OK")
	})
}
