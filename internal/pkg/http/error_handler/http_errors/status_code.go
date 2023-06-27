package http_errors

import (
	"errors"
	"fmt"
	"net/http"
	"server-template/internal/models/err_const"
	"strings"
)

func statusCodeAndErrorMessage(err error) (int, string, string) {
	switch {
	case errors.Is(err, err_const.ErrJsonUnMarshal):
		return http.StatusBadRequest, err_const.ErrJsonUnMarshal.Error(), err.Error()
	case errors.Is(err, err_const.ErrJsonMarshal):
		return http.StatusInternalServerError, err_const.ErrJsonMarshal.Error(), err.Error()
	case errors.Is(err, err_const.ErrDatabaseRecordNotFound):
		return http.StatusNotFound, err_const.ErrDatabaseRecordNotFound.Error(), err.Error()
	case errors.Is(err, err_const.ErrMissingUser):
		return http.StatusBadRequest, err_const.ErrMissingUser.Error(), err.Error()
	case errors.Is(err, err_const.ErrMissingToken):
		return http.StatusBadRequest, err_const.ErrMissingToken.Error(), err.Error()
	case strings.HasSuffix(err.Error(), "record not found"):
		return http.StatusNotFound, "запись не найдена", err.Error()
	default:
		return http.StatusInternalServerError, fmt.Errorf("ошибка обработки запроса: %w", err).Error(), err.Error()
	}
}
