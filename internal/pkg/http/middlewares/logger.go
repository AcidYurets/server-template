package middlewares

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"server-template/internal/modules/logger"
	"time"

	"go.uber.org/zap"
)

func UseContextLogger(lg *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := logger.SetToCtx(c.Request().Context(), lg)
			c.Request().WithContext(ctx)

			return next(c)
		}
	}
}

func UseRequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var start, stop time.Time

			start = time.Now()
			err := next(c)
			stop = time.Now()

			// Обрабатываем ошибку
			if err != nil {
				// Вынимаем из echo.HTTPError обернутую ошибку, если она там указана
				var echoError *echo.HTTPError
				if errors.As(err, &echoError) {
					wrappedErr := echoError.Unwrap()
					if wrappedErr != nil {
						err = wrappedErr
					}
				}

				// Вызываем установленный обработчик
				c.Error(err)
			}

			fields := []zap.Field{
				zap.String("ip", c.RealIP()),
				zap.String("path", c.Path()),
				zap.String("method", c.Request().Method),
				zap.Int("status", c.Response().Status),
				zap.String("latency", stop.Sub(start).String()),
			}

			body, readErr := io.ReadAll(c.Request().Body)
			if readErr != nil {
				fields = append(fields, zap.NamedError("read_body_error", fmt.Errorf("не удалось прочитать тело запроса: %w", err)))
			} else {
				fields = append(fields, zap.ByteString("body", body))
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
			}

			s := c.Response().Status
			switch {
			case s >= 500:
				msg := fmt.Sprintf("Неизвестная внутренняя ошибка")
				if err != nil {
					msg = fmt.Sprintf("Внутренняя ошибка сервера: %s", err.Error())
				}
				logger.Error(msg, fields...)
			case s >= 400:
				msg := fmt.Sprintf("Неизвестная ошибка в запросе")
				if err != nil {
					msg = fmt.Sprintf("Ошибка в запросе: %s", err.Error())
				}
				logger.Warn(msg, fields...)
			default:
				logger.Info("Запрос выполнен успешно", fields...)
			}

			return nil
		}
	}
}
