package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"nnw_s/config"
)

type UserRegistrationDataRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyRegistrationEmailCode struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type VerifyRegistrationCodeRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
	Code     string `json:"code" validate:"required"`
	Uid      string `json:"uid" validate:"required"`
}

type UserLoginDataRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
}

type VerifyLoginCodeRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
	Code     string `json:"code" validate:"required"`
}

type CheckEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

func (h *Handler) verifyRegistrationEmail(c echo.Context) error {
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
	user, _ := h.services.GetUserByEmail(registrationUserData.Email)
	if user != nil {
		return c.JSON(http.StatusBadRequest, UserAlreadyExist)
	}

	err = h.services.CreateEmail(registrationUserData.Email, "verify")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.NoContent(200)
}

func (h *Handler) verifyRegistrationEmailResend(c echo.Context) error {
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

	err = h.services.CreateEmail(registrationUserData.Email, "verify")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.NoContent(200)
}

func (h *Handler) verifyRegistrationEmailCode(c echo.Context) error {
	var verifyRegistrationEmailCode VerifyRegistrationEmailCode

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyRegistrationEmailCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(verifyRegistrationEmailCode)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	valid, err := h.services.CheckEmailCode(verifyRegistrationEmailCode.Email, verifyRegistrationEmailCode.Code, "verify")
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Generate 2FA Image
	buffImg, key, err := h.services.Generate2FaImage(verifyRegistrationEmailCode.Email)
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
	valid := h.services.Check2FaCode(verifyRegistrationCode.Code, templateData.TwoFAS)
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Create User
	_, err = h.services.CreateUser(verifyRegistrationCode.Email, verifyRegistrationCode.Password, templateData.TwoFAS)
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
	user, err := h.services.GetUserByEmail(userLoginData.Email)
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
	user, err := h.services.GetUserByEmail(verifyLoginCode.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, UserNotFound)
	}

	// Check Valid 2FA Code
	valid := h.services.Check2FaCode(verifyLoginCode.Code, user.SecretOTPKey)
	if !valid {
		return c.JSON(http.StatusBadRequest, InvalidCode)
	}

	// Create JWT
	jwtToken, err := h.services.CreateJWTToken(verifyLoginCode.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": jwtToken})
}

func (h *Handler) checkEmail(c echo.Context) error {
	var checkEmailData CheckEmailRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkEmailData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidJson)
	}

	// Translation
	trans := config.ValidatorConfig(h.validate)

	// Validate Body
	err = h.validate.Struct(checkEmailData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, invalidValidationFieldsArray(errArray))
	}

	// Find User
	user, err := h.services.GetUserByEmail(checkEmailData.Email)
	if user == nil {
		return c.NoContent(http.StatusOK)
	}

	if user.Email == checkEmailData.Email {
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
