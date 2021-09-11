package auth

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"nnw_s/config"
	"nnw_s/internal/user"
	"nnw_s/pkg/errors"
	"strings"
	"time"
)

type UserRegistrationDataRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
}

type SendVerifyRegistrationEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResendEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckRegistrationEmailCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type Generate2FaCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckRegistrationCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type FinishRegistrationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserLoginDataRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwd"`
}

type CheckLogin2FACodeRequest struct {
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

func (h *Handler) preRegistration(c echo.Context) error {
	var registrationUserData UserRegistrationDataRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&registrationUserData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(registrationUserData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find active user
	activeUser, _ := h.userService.GetActiveUser(context.Background(), registrationUserData.Email)
	if activeUser != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrAlreadyExists, "User already exist!"))
	}

	// Find disable user
	disableUser, _ := h.userService.GetDisableUser(context.Background(), registrationUserData.Email)
	if disableUser != nil {
		disableUser.Email = registrationUserData.Email
		disableUser.Password = registrationUserData.Password

		err = h.userService.UpdateDisableUser(context.Background(), disableUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, "Invalid user data!"))
		}
	} else {
		var userDto user.CreateUserDTO
		userDto.Email = registrationUserData.Email
		userDto.Password = registrationUserData.Password
		userDto.SecretOTP = "null"

		_, err := h.userService.CreateUser(context.Background(), &userDto)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, "Invalid user data!"))
		}
	}

	return c.NoContent(200)
}

func (h *Handler) sendVerifyRegistrationEmail(c echo.Context) error {
	var verifyEmailData SendVerifyRegistrationEmailRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&verifyEmailData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(verifyEmailData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find disable user
	disableUser, _ := h.userService.GetDisableUser(context.Background(), verifyEmailData.Email)
	if disableUser != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrNotFound, "User not found!"))
	}

	err = h.authService.CreateEmail(context.Background(), verifyEmailData.Email, "verify")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(200)
}

func (h *Handler) resendVerifyRegistrationEmail(c echo.Context) error {
	var resendEmailData ResendEmailRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&resendEmailData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(resendEmailData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	err = h.authService.CreateEmail(context.Background(), resendEmailData.Email, "verify")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(200)
}

func (h *Handler) checkRegistrationEmailCode(c echo.Context) error {
	var checkRegistrationEmailCodeData CheckRegistrationEmailCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkRegistrationEmailCodeData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(checkRegistrationEmailCodeData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	valid, err := h.authService.CheckEmailCode(context.Background(), checkRegistrationEmailCodeData.Email, checkRegistrationEmailCodeData.Code, "verify")
	if !valid {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidCode, "Invalid email code!"))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) generate2FAQrCode(c echo.Context) error {
	var generate2FAQrCodeData Generate2FaCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&generate2FAQrCodeData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(generate2FAQrCodeData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find disable user
	disableUser, _ := h.userService.GetDisableUser(context.Background(), generate2FAQrCodeData.Email)
	if disableUser != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrNotFound, "User not found!"))
	}

	// Generate 2FA Image
	buffImg, key, err := h.authService.Generate2FaImage(context.Background(), generate2FAQrCodeData.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	// Update user SecretOTP key
	disableUser.SecretOTPKey = key.Secret()
	disableUser.UpdatedAt = time.Now()
	err = h.userService.UpdateUser(context.Background(), disableUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	// Set Header
	c.Response().Header().Set("Content-Type", "image/png")
	//c.Response().Header().Set("Access-Control-Expose-Headers", "Tid")

	// Write image bytes
	_, err = c.Response().Write(buffImg.Bytes())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewInternal("Internal server error!"))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) checkRegistration2FaCode(c echo.Context) error {
	var checkRegistrationCodeData CheckRegistrationCodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkRegistrationCodeData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(checkRegistrationCodeData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find disable user
	disableUser, _ := h.userService.GetDisableUser(context.Background(), checkRegistrationCodeData.Email)
	if disableUser != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrNotFound, "User not found!"))
	}

	// Check Valid 2FA Code
	valid := h.authService.Check2FaCode(checkRegistrationCodeData.Code, disableUser.SecretOTPKey)
	if !valid {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidCode, "Invalid 2FA code!"))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) finishRegistration(c echo.Context) error {
	var finishRegistrationData FinishRegistrationRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&finishRegistrationData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(finishRegistrationData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find disable user
	disableUser, _ := h.userService.GetDisableUser(context.Background(), finishRegistrationData.Email)
	if disableUser != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrNotFound, "User not found!"))
	}

	// Create User
	disableUser.Status = "active"
	err = h.userService.UpdateUser(context.Background(), disableUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, "Invalid user data!"))
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) login(c echo.Context) error {
	var userLoginData UserLoginDataRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&userLoginData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(userLoginData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find user
	userByEmail, err := h.userService.GetUserByEmail(context.Background(), userLoginData.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrNotFound, "User not found!"))
	}

	// Check password
	err = h.authService.CheckPassword(context.Background(), userLoginData.Password, userByEmail.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, "Incorrect login data!"))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) checkLogin2faCode(c echo.Context) error {
	var checkLogin2FACodeData CheckLogin2FACodeRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkLogin2FACodeData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(checkLogin2FACodeData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find user
	userByEmail, err := h.userService.GetUserByEmail(context.Background(), checkLogin2FACodeData.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrNotFound, "User not found!"))
	}

	// Check Valid 2FA Code
	valid := h.authService.Check2FaCode(checkLogin2FACodeData.Code, userByEmail.SecretOTPKey)
	if !valid {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidCode, "Invalid 2FA code!"))
	}

	// Create JWT
	jwtToken, err := h.authService.CreateJWTToken(context.Background(), checkLogin2FACodeData.Email)
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
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(checkEmailData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	// Find User
	userByEmail, err := h.userService.GetUserByEmail(context.Background(), checkEmailData.Email)
	if userByEmail == nil {
		return c.NoContent(http.StatusOK)
	}

	if userByEmail.Email == checkEmailData.Email {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func (h *Handler) checkJwt(c echo.Context) error {
	var checkTokenData CheckTokenRequest

	// Parse User Data
	err := json.NewDecoder(c.Request().Body).Decode(&checkTokenData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidJson, "Invalid json!"))
	}

	// Translation
	trans := config.ValidatorConfig(h.validate, h.cfg)

	// Validate Body
	err = h.validate.Struct(checkTokenData)
	if err != nil {
		var errArray []string
		for _, e := range err.(validator.ValidationErrors) {
			errArray = append(errArray, e.Translate(trans))
		}

		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidData, strings.Join(errArray, ",")))
	}

	_, err = h.authService.VerifyJWTToken(context.Background(), checkTokenData.Token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.WithMessage(ErrWrongToken, "Wrong jwt token!"))
	}

	return c.NoContent(http.StatusOK)
}
