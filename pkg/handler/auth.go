package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

func (h *Handler) registration(c echo.Context) error {
	userData := new(userData)

	// Parse User Data
	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Find user
	user, _ := h.services.GetUserByLogin(userData.Login)
	if user != nil {
		return c.JSON(http.StatusBadRequest, UserAlreadyExist)
	}

	// Generate 2FA Image
	buffImg, key, err := h.services.Generate2FaImage(userData.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	// Save template uid with secret key
	templateId, err := h.services.CreateTemplateUserData(key.Secret())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	// Set Header
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Access-Control-Expose-Headers", "Tid")
	c.Response().Header().Set("Tid", *templateId)

	// Write image bytes
	_, err = c.Response().Write(buffImg.Bytes())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) verifyRegistration2FaCode(c echo.Context) error {
	userData := new(userData)

	// Parse User Data
	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Get Secret
	templateData, err := h.services.GetTemplateUserDataById(userData.Uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	// Check Valid 2FA Code
	valid := totp.Validate(userData.Code, templateData.TwoFAS)
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Create User
	_, err = h.services.CreateUser(userData.Login, userData.Email, userData.Password, templateData.TwoFAS)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) login(c echo.Context) error {
	userData := new(userData)

	// Parse User Data
	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Find user
	user, err := h.services.GetUserByLogin(userData.Login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, UserNotFound)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidPassword)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) verifyLogin2fa(c echo.Context) error {
	userData := new(userData)

	// Parse User Data
	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Find user
	user, err := h.services.GetUserByLogin(userData.Login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, UserNotFound)
	}

	// Check Valid 2FA Code
	valid := totp.Validate(userData.Code, user.SecretOTPKey)
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Create JWT
	jwtToken, err := h.services.CreateJWTToken(userData.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": jwtToken})
}

func (h *Handler) checkUserName(c echo.Context) error {
	userData := new(userData)

	// Parse User Data
	if err := c.Bind(userData); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Find User
	user, _ := h.services.GetUserByLogin(userData.Login)
	if user == nil {
		return c.NoContent(http.StatusOK)
	}

	if user.Login == userData.Login {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func (h *Handler) checkJwt(c echo.Context) error {
	token := new(Token)

	// Parse User Data
	if err := c.Bind(token); err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	_, err := h.services.VerifyJWTToken(token.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, WrongToken)
	}

	return c.NoContent(http.StatusOK)
}
