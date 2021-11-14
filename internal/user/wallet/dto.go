package wallet

import (
	"github.com/go-playground/validator/v10"
	"math/big"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/helpers"
	"time"
)

const passwordMinLength = 8

func Validate(dto interface{}, shift int) error {
	validate := validator.New()

	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		decodedPassword, err := helpers.CaesarShift(password, -shift)
		if err != nil {
			return false
		}

		if len(decodedPassword) < passwordMinLength {
			return false
		}

		return true
	})

	if err := validate.Struct(dto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.WithMessage(ErrInvalidRequest, err.Error())
		}

		validationErr := ErrInvalidRequest
		for _, err := range err.(validator.ValidationErrors) {
			validationErr = errors.WithMessage(validationErr, err.Error())
		}
		return validationErr
	}
	return nil
}

type BalanceDTO struct {
	Balance    float64  `json:"balance"`
	BalanceInt *big.Int `json:"balance_int"`
	BalanceStr string   `json:"balance_str"`
}

type TxsDTO struct {
	//ToAddress   string  `json:"to_address"`
	//FromAddress string  `json:"from_address"`
	//Amount      float64 `json:"amount"`
	//TxId        string  `json:"tx_id"`
	//Address  string      `json:"address"`
	//Category string      `json:"category"`
	//Amount   interface{} `json:"amount"`
	Txid   string        `json:"txid"`
	Input  []*InputTxDTO `json:"input"`
	Output []*OutTxDTO   `json:"output"`
	Time   time.Time     `json:"time"`
}

type InputTxDTO struct {
	//Vin     int64   `json:"vin"`
	//TxId    string  `json:"txid"`
	Address string  `json:"address"`
	Value   float64 `json:"value"`
}

type OutTxDTO struct {
	Address string  `json:"address"`
	Value   float64 `json:"value"`
}

type CreateWalletDTO struct {
	Password string `json:"password" validate:"required,password"`
	Backup   *bool  `json:"backup" validate:"required"`
	Jwt      string `json:"jwt" validate:"required"`
}

type GetWalletDTO struct {
	Jwt      string `json:"jwt" validate:"required"`
	WalletId string `json:"wallet_id" validate:"required"`
}

type UnlockWalletDTO struct {
	Jwt      string `json:"jwt" validate:"required"`
	Name     string `json:"name" validate:"required"`
	WalletId string `json:"wallet_id" validate:"required"`
}

type GetWalletBalanceDTO struct {
	Jwt      string `json:"jwt" validate:"required"`
	Name     string `json:"name" validate:"required"`
	WalletId string `json:"wallet_id" validate:"required"`
}

type GetWalletTxDTO struct {
	Jwt      string `json:"jwt" validate:"required"`
	Name     string `json:"name" validate:"required"`
	WalletId string `json:"wallet_id" validate:"required"`
	Address  string `json:"address" validate:"required"`
}
