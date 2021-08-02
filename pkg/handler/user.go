package handler

import (
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"image/png"
)

type userData struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Uid      string `json:"uid"`
}

func (h *Handler) registration(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Generate 2FA Image
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "NNW",
		AccountName: "example@examole.com",
	})

	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Encode image
	err = png.Encode(&buf, img)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Save template uid with secret key
	templateId, err := h.services.CreateTemplateUserData(key.Secret())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Set Header
	ctx.Response().Header.Set("Content-Type", "image/png")
	ctx.Response().Header.Set("Access-Control-Expose-Headers", "Tid")
	ctx.Response().Header.Set("Tid", *templateId)
	_, err = ctx.Write(buf.Bytes())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) verify2FaCode(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Get Secret
	templateData, err := h.services.GetTemplateUserDataById(userData.Uid)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Check Valid 2FA Code
	valid := totp.Validate(userData.Code, templateData.TwoFAS)
	if !valid {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{"error": errors.New(" Invalid code!").Error()})
	}

	// Create User
	_, err = h.services.CreateUser(userData.Login, userData.Email, userData.Password, templateData.TwoFAS)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{"error": errors.New(" Invalid code!").Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
