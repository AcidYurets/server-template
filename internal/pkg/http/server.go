package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

type Config struct {
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	DevelopMode  bool

	ErrorHandler       echo.HTTPErrorHandler
	MiddlewareHandlers []echo.MiddlewareFunc
}

func New(cfg Config) *echo.Echo {
	e := echo.New()

	e.Use(middleware.BodyLimit(fmt.Sprintf("%dM", 128)))
	e.Server.ReadTimeout = cfg.ReadTimeOut
	e.Server.WriteTimeout = cfg.WriteTimeOut
	e.Server.IdleTimeout = cfg.ReadTimeOut
	e.HTTPErrorHandler = cfg.ErrorHandler

	e.Use(cfg.MiddlewareHandlers...)

	return e
}
