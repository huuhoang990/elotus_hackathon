package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ApiResponse map[string]any

type JSONSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONFailedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendFailedValidationResponse(c echo.Context, message string, errors interface{}) error {
	return c.JSON(http.StatusBadRequest, JSONFailedResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
