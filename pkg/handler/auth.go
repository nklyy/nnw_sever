package handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type userData struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Uid      string `json:"uid"`
}

type Token struct {
	Token string `json:"token"`
}

func (h *Handler) registration(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Generate 2FA Image
	buffImg, key, err := h.services.Generate2FaImage(userData.Login)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Save template uid with secret key
	templateId, err := h.services.CreateTemplateUserData(key.Secret())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Set Header
	ctx.Response().Header.Set("Content-Type", "image/png")
	ctx.Response().Header.Set("Access-Control-Expose-Headers", "Tid")
	ctx.Response().Header.Set("Tid", *templateId)

	// Write image bytes
	_, err = ctx.Write(buffImg.Bytes())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) verifyRegistration2FaCode(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Get Secret
	templateData, err := h.services.GetTemplateUserDataById(userData.Uid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Check Valid 2FA Code
	valid := totp.Validate(userData.Code, templateData.TwoFAS)
	if !valid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Invalid code!").Error()})
	}

	// Create User
	_, err = h.services.CreateUser(userData.Login, userData.Email, userData.Password, templateData.TwoFAS)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Invalid code!").Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) login(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Find user
	user, err := h.services.GetUserByLogin(userData.Login)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) verifyLogin2fa(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Find user
	user, err := h.services.GetUserByLogin(userData.Login)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	// Check Valid 2FA Code
	valid := totp.Validate(userData.Code, user.SecretOTPKey)
	if !valid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Invalid code!").Error()})
	}

	// Create JWT
	jwtToken, err := h.services.CreateJWTToken(userData.Login)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Something wrong!").Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"token": jwtToken})
}

func (h *Handler) checkLogin(ctx *fiber.Ctx) error {
	var userData userData

	// Parse User Data
	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	// Find User
	user, _ := h.services.GetUserByLogin(userData.Login)
	if user == nil {
		return ctx.SendStatus(fiber.StatusOK)
	}

	if user.Login == userData.Login {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusInternalServerError)
}

func (h *Handler) checkJwt(ctx *fiber.Ctx) error {
	var token Token

	// Parse User Data
	if err := ctx.BodyParser(&token); err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errors.New(" Invalid json!").Error()})
	}

	_, err := h.services.VerifyJWTToken(token.Token)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors.New(" Wrong token!").Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
