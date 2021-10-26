package wallet

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"nnw_s/internal/auth/jwt"
	"nnw_s/pkg/errors"
)

type Handler struct {
	walletSvc Service
	jwtSvc    jwt.Service
	shift     int
}

func NewHandler(walletSvc Service, jwtSvc jwt.Service, shift int) *Handler {
	return &Handler{
		walletSvc: walletSvc,
		jwtSvc:    jwtSvc,
		shift:     shift,
	}
}

func (h *Handler) SetupRoutes(router *echo.Echo) {
	v1 := router.Group("/api/v1")

	// Create wallet
	v1.POST("/create-wallet", h.createWallet)
}

func (h *Handler) createWallet(ctx echo.Context) error {
	var dto CreateWalletDTO

	if err := ctx.Bind(&dto); err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	jwtPayload, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	fmt.Println(jwtPayload)
	fmt.Println(*dto.Backup)

	walletPayload, err := h.walletSvc.CreateWallet(ctx.Request().Context(), &dto, jwtPayload.Email, h.shift)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	fmt.Println(walletPayload)

	return ctx.NoContent(200)
}
