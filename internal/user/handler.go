package user

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"nnw_s/internal/auth/jwt"
	"nnw_s/pkg/errors"
)

type Handler struct {
	userSvc Service
	jwtSvc  jwt.Service
	shift   int
}

func NewHandler(userSvc Service, jwtSvc jwt.Service, shift int) *Handler {
	return &Handler{
		userSvc: userSvc,
		jwtSvc:  jwtSvc,
		shift:   shift,
	}
}

func (h *Handler) SetupRoutes(router *echo.Echo) {
	v1 := router.Group("/api/v1")

	// Create wallet
	v1.POST("/get-user", h.getUser)
}

func (h *Handler) getUser(ctx echo.Context) error {
	var dto GetUserDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	jwtPayload, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	user, err := h.userSvc.GetUserByEmail(ctx.Request().Context(), jwtPayload.Email)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.JSON(http.StatusOK, NormalizeGetUserResponseDTO(user))
}
