package wallet

import (
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

	// Get wallet
	v1.POST("/get-wallet", h.getWallet)

	// Create wallet
	v1.POST("/create-wallet", h.createWallet)

	//Wallet
	v1.POST("/get-balance", h.getBalance)
	//v1.POST("/unlock", h.unlockWallet)
	//v1.POST("/get-tx", h.getTX)
}

func (h *Handler) createWallet(ctx echo.Context) error {
	var dto CreateWalletDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	jwtPayload, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	walletPayload, err := h.walletSvc.CreateWallet(ctx.Request().Context(), &dto, jwtPayload.Email, h.shift)
	if err != nil {
		return ctx.JSON(errors.HTTPCode(err), err)
	}

	return ctx.JSON(http.StatusOK, walletPayload)
}

func (h *Handler) getWallet(ctx echo.Context) error {
	var dto GetWalletDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	jwtPayload, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	walletPayload, err := h.walletSvc.GetWallet(ctx.Request().Context(), jwtPayload.Email, dto.WalletId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, walletPayload)
}

//func (h *Handler) unlockWallet(ctx echo.Context) error {
//	var dto UnlockWalletDTO
//
//	if err := ctx.Bind(&dto); err != nil {
//		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
//	}
//
//	if err := Validate(dto, h.shift); err != nil {
//		return ctx.JSON(http.StatusBadRequest, err)
//	}
//
//	jwtPayload, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
//	if err != nil {
//		return ctx.JSON(http.StatusBadRequest, err)
//	}
//	switch dto.Name {
//	case "BTC":
//		// Call get unlock wallet from btc rpc
//		fmt.Println("BTC")
//	}
//
//}

func (h *Handler) getBalance(ctx echo.Context) error {
	var dto GetWalletBalanceDTO

	if err := ctx.Bind(&dto); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
	}

	if err := Validate(dto, h.shift); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	_, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	balance, err := h.walletSvc.GetBalance(ctx.Request().Context(), &dto)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, balance)
}

//func (h *Handler) getTX(ctx echo.Context) error {
//	var dto GetWalletTxDTO
//
//	if err := ctx.Bind(&dto); err != nil {
//		return ctx.JSON(http.StatusBadRequest, errors.WithMessage(ErrInvalidRequest, err.Error()))
//	}
//
//	if err := Validate(dto, h.shift); err != nil {
//		return ctx.JSON(http.StatusBadRequest, err)
//	}
//
//	jwtPayload, err := h.jwtSvc.VerifyJWT(ctx.Request().Context(), dto.Jwt)
//	if err != nil {
//		return ctx.JSON(http.StatusBadRequest, err)
//	}
//
//	switch dto.Name {
//	case "BTC":
//		// Call get transaction list tx from btc rpc
//		fmt.Println("BTC")
//	}
//}