package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	InternalServerError = echo.Map{"error": errors.New(" Something wrong!").Error()}

	BadRequest = echo.Map{"error": errors.New(" Bad Request!").Error()}

	InvalidJson     = echo.Map{"error": errors.New(" Invalid json!").Error()}
	InvalidCode     = echo.Map{"error": errors.New(" Invalid code!").Error()}
	InvalidPassword = echo.Map{"error": errors.New(" Invalid password!").Error()}

	UserAlreadyExist = echo.Map{"error": errors.New(" User already exist!").Error()}
	UserNotFound     = echo.Map{"error": errors.New(" User not found!").Error()}

	WrongToken = echo.Map{"error": errors.New(" Wrong token!").Error()}
)
