package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"nnw_s/config"
	"nnw_s/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	// Find user
	user, _ := h.userService.GetUserByEmail(context.Background(), registrationUserData.Email)
	if user != nil {
		return c.JSON(http.StatusBadRequest, common.UserAlreadyExist)
	}

	err = h.authService.CreateEmail(context.Background(), registrationUserData.Email, "verify")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(200)
}

func (h *Handler) verifyRegistrationEmailResend(c echo.Context) error {
	var registrationUserData UserRegistrationDataRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&registrationUserData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	err = h.authService.CreateEmail(context.Background(), registrationUserData.Email, "verify")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(200)
}

func (h *Handler) verifyRegistrationEmailCode(c echo.Context) error {
	var verifyRegistrationEmailCode VerifyRegistrationEmailCode

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyRegistrationEmailCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	valid, err := h.authService.CheckEmailCode(context.Background(), verifyRegistrationEmailCode.Email, verifyRegistrationEmailCode.Code, "verify")
	if !valid {
		return c.JSON(http.StatusBadRequest, common.InvalidCode)
	}

	// Generate 2FA Image
	buffImg, key, err := h.authService.Generate2FaImage(context.Background(), verifyRegistrationEmailCode.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	// Save template uid with secret key
	templateId, err := h.authService.CreateTemplateUserData(context.Background(), key.Secret())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	// Set Header
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Access-Control-Expose-Headers", "Tid")
	c.Response().Header().Set("Tid", templateId)

	// Write image bytes
	_, err = c.Response().Write(buffImg.Bytes())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) verifyRegistration2FaCode(c echo.Context) error {
	var verifyRegistrationCode VerifyRegistrationCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyRegistrationCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	// Get Secret
	templateData, err := h.userService.GetTemplateUserDataByID(context.Background(), verifyRegistrationCode.Uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	// Check Valid 2FA Code
	valid := h.authService.Check2FaCode(verifyRegistrationCode.Code, templateData.TwoFAS)
	if !valid {
		return c.JSON(http.StatusBadRequest, common.InvalidCode)
	}

	// Create User
	_, err = h.userService.CreateUser(context.Background(), verifyRegistrationCode.Email, verifyRegistrationCode.Password, templateData.TwoFAS)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidData)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) login(c echo.Context) error {
	var userLoginData UserLoginDataRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&userLoginData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	// Find user
	user, err := h.userService.GetUserByEmail(context.Background(), userLoginData.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.UserNotFound)
	}

	// Check password
	err = h.authService.CheckPassword(context.Background(), userLoginData.Password, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) verifyLogin2fa(c echo.Context) error {
	var verifyLoginCode VerifyLoginCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyLoginCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	// Find user
	user, err := h.userService.GetUserByEmail(context.Background(), verifyLoginCode.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.UserNotFound)
	}

	// Check Valid 2FA Code
	valid := h.authService.Check2FaCode(verifyLoginCode.Code, user.SecretOTPKey)
	if !valid {
		return c.JSON(http.StatusBadRequest, common.InvalidCode)
	}

	// Create JWT
	jwtToken, err := h.authService.CreateJWTToken(context.Background(), verifyLoginCode.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.JSON(http.StatusOK, echo.Map{"token": jwtToken})
}

func (h *Handler) checkEmail(c echo.Context) error {
	var checkEmailData CheckEmailRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkEmailData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	// Find User
	user, err := h.userService.GetUserByEmail(context.Background(), checkEmailData.Email)
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
		return c.JSON(http.StatusBadRequest, common.InvalidJson)
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

		return c.JSON(http.StatusBadRequest, common.InvalidValidationFieldsArray(errArray))
	}

	_, err = h.authService.VerifyJWTToken(context.Background(), checkTokenData.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.WrongToken)
	}

	return c.NoContent(http.StatusOK)
}
