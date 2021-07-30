package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type userData struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) registration(ctx *fiber.Ctx) error {
	var userData userData

	if err := ctx.BodyParser(&userData); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	fmt.Println(userData)

	return ctx.SendStatus(fiber.StatusCreated)
}
