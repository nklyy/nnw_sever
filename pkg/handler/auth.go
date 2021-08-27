package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"net/http"
	"nnw_s/config"
)

type UserRegistrationDataRequest struct {
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
}

type VerifyRegistrationCodeRequest struct {
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
	Code     string `json:"code" validate:"required"`
	Uid      string `json:"uid" validate:"required"`
}

type UserLoginDataRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,passwd"`
}

type VerifyLoginCodeRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,passwd"`
	Code     string `json:"code" validate:"required"`
}

type CheckLoginRequest struct {
	Login string `json:"login" validate:"required"`
}

type CheckTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

func (h *Handler) registration(c echo.Context) error {
	var registrationUserData UserRegistrationDataRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&registrationUserData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(registrationUserData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	// Find user
	user, _ := h.services.GetUserByLogin(registrationUserData.Login)
	if user != nil {
		return c.JSON(http.StatusBadRequest, UserAlreadyExist)
	}

	// Generate 2FA Image
	buffImg, key, err := h.services.Generate2FaImage(registrationUserData.Login)
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
	var verifyRegistrationCode VerifyRegistrationCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyRegistrationCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(verifyRegistrationCode)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	// Get Secret
	templateData, err := h.services.GetTemplateUserDataById(verifyRegistrationCode.Uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	// Check Valid 2FA Code
	valid := totp.Validate(verifyRegistrationCode.Code, templateData.TwoFAS)
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Create User
	_, err = h.services.CreateUser(verifyRegistrationCode.Login, verifyRegistrationCode.Email, verifyRegistrationCode.Password, templateData.TwoFAS)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidData)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) login(c echo.Context) error {
	var userLoginData UserLoginDataRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&userLoginData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(userLoginData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	// Find user
	user, err := h.services.GetUserByLogin(userLoginData.Login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, UserNotFound)
	}

	// Check password
	validPass, err := h.services.CheckPassword(userLoginData.Password, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	if !validPass {
		return c.JSON(http.StatusBadRequest, InvalidPassword)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) verifyLogin2fa(c echo.Context) error {
	var verifyLoginCode VerifyLoginCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyLoginCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(verifyLoginCode)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	// Find user
	user, err := h.services.GetUserByLogin(verifyLoginCode.Login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, UserNotFound)
	}

	// Check Valid 2FA Code
	valid := totp.Validate(verifyLoginCode.Code, user.SecretOTPKey)
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Create JWT
	jwtToken, err := h.services.CreateJWTToken(verifyLoginCode.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": jwtToken})
}

func (h *Handler) checkLogin(c echo.Context) error {
	var checkLoginData CheckLoginRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkLoginData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(checkLoginData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	// Find User
	user, err := h.services.GetUserByLogin(checkLoginData.Login)
	if user == nil {
		return c.NoContent(http.StatusOK)
	}

	if user.Login == checkLoginData.Login {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func (h *Handler) checkJwt(c echo.Context) error {
	var checkTokenData CheckTokenRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkTokenData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(checkTokenData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	_, err = h.services.VerifyJWTToken(checkTokenData.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, WrongToken)
	}

	return c.NoContent(http.StatusOK)
}
