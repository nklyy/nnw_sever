package auth

import (
	"net/http"
	"nnw_s/internal/auth/jwt"
	"nnw_s/pkg/errors"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	registrationSvc RegistrationService
	loginSvc        LoginService
	jwtSvc          jwt.Service
	shift           int
}

func NewHandler(registrationSvc RegistrationService, loginSvc LoginService, jwtSvc jwt.Service, shift int) *Handler {
	return &Handler{
		registrationSvc: registrationSvc,
		loginSvc:        loginSvc,
		jwtSvc:          jwtSvc,
		shift:           shift,
	}
}

func (h *Handler) SetupRoutes(router *echo.Echo) {
	v1 := router.Group("/api/v1")

	// Registration and Verify Email
	v1.POST("/register-user", h.registerUser)
	v1.POST("/verify-user", h.verifyUser)
	v1.POST("/resend-verification-email", h.resendVerificationRegistrationEmail)
	v1.POST("/setup-twoFa", h.setupTwoFA)
	v1.POST("/activate-user", h.activateUser)

	// Login and Logout
	v1.POST("/login", h.login)
	v1.POST("/login-code", h.loginCode)
	v1.POST("/logout", h.logout)

	// Validate JWT Token
	v1.POST("/validate-token", h.validateToken)

	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}

func (h *Handler) registerUser(ctx echo.Context) error {
	var dto RegisterUserDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	if err := h.registrationSvc.RegisterUser(ctx.Request().Context(), &dto); err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.NoContent(200)
}

func (h *Handler) verifyUser(ctx echo.Context) error {
	var dto VerifyUserDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := h.registrationSvc.VerifyUser(ctx.Request().Context(), &dto); err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.NoContent(200)
}

func (h *Handler) resendVerificationRegistrationEmail(ctx echo.Context) error {
	var dto ResendActivationEmailDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := h.registrationSvc.ResendVerificationEmail(ctx.Request().Context(), &dto); err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.NoContent(200)
}

func (h *Handler) setupTwoFA(ctx echo.Context) error {
	var dto SetupTwoFaDTO

	if err := ctx.Bind(&dto); err != nil {

		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	buf, err := h.registrationSvc.SetupTwoFA(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	// write image bytes
	ctx.Response().Header().Set("Content-Type", "image/png")
	if _, err = ctx.Response().Write(buf); err != nil {
		return ctx.JSON(http.StatusInternalServerError, errors.NewInternal(err.Error()))
	}
	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) activateUser(ctx echo.Context) error {
	var dto ActivateUserDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := h.registrationSvc.ActivateUser(ctx.Request().Context(), &dto); err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (h *Handler) login(ctx echo.Context) error {
	var dto LoginDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err := h.loginSvc.Login(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Handler) loginCode(ctx echo.Context) error {
	var dto LoginCodeDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	tokenDTO, err := h.loginSvc.CheckCode(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.JSON(http.StatusOK, tokenDTO)
}

func (h *Handler) validateToken(ctx echo.Context) error {
	var dto ValidateTokenDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	_, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Token)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.NoContent(http.StatusOK)
}

// todo: implement
func (h *Handler) logout(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
