package server_time

import (
	"github.com/labstack/echo/v4"
	"time"
)

func SetDateHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			c.Response().Header().Set("Date", time.Now().UTC().Format(time.RFC1123))
		}()
		return next(c)
	}
}
