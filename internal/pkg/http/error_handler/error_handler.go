package error_handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"server-template/internal/pkg/http/error_handler/http_errors"
)

func ErrorHandler(err error, c echo.Context) {
	var (
		apiError  *http_errors.ErrorResponse
		echoError *echo.HTTPError // Ошибка обработки роутинга и других ошибок
	)

	switch {
	case errors.As(err, &echoError):
		apiError = &http_errors.ErrorResponse{
			Code:    echoError.Code,
			Message: "внутренняя ошибка сервера",
			Details: echoError.Error(),
		}
	case errors.As(err, &apiError):
		if apiError.Code == 0 {
			apiError.Code = http.StatusInternalServerError
		}
	default:
		apiError = http_errors.NewErrorResponse(err)
	}

	// TODO: Обработать ошибку
	_ = c.JSON(apiError.Code, apiError)
}
