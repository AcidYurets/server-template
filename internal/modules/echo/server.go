package echo

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"server-template/internal/modules/config"
	"server-template/internal/modules/logger"
	"server-template/internal/pkg/http"
	"server-template/internal/pkg/http/error_handler"
	"server-template/internal/pkg/http/middlewares"
	"server-template/internal/pkg/http/middlewares/server_time"
)

func NewEcho(cfg config.Config, logger *zap.Logger) *echo.Echo {
	echoServer := newHTTPServer(cfg, logger)

	return echoServer
}

func InvokeEcho(
	app *echo.Echo,
	cfg config.Config,
	lifecycle fx.Lifecycle,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go runServer(app, cfg)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown(ctx)
		},
	})
}

func runServer(app *echo.Echo, cfg config.Config) error {
	address := getAddress(cfg)
	return app.Start(address)
}

// getAddress метод получения адреса для сервера
func getAddress(config config.Config) string {
	return fmt.Sprintf("%s:%s", config.HTTPServerHost, config.HTTPServerPort)
}

func newHTTPServer(config config.Config, lg *zap.Logger) *echo.Echo {
	settings := config

	uses := []echo.MiddlewareFunc{
		server_time.SetDateHeader,
		middlewares.UseRequestLogger(lg),
		middlewares.UseContextLogger(lg),
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				logger.GetFromCtx(c.Request().Context()).
					With(zap.String("panic", error.Error(err)), zap.ByteString("stacktrace", stack)).
					Error("паника при обработке запроса")

				return err
			},
		}),

		// Разрешение на подключение со всех адресов
		// Отключение CORS
		middleware.CORS(),
	}

	echoServer := http.New(http.Config{
		ReadTimeOut:        settings.HTTPServerReadTimeOut,
		WriteTimeOut:       settings.HTTPServerWriteTimeOut,
		DevelopMode:        settings.DevMode,
		ErrorHandler:       error_handler.ErrorHandler,
		MiddlewareHandlers: uses,
	})

	echoServer.Any("/ntp", func(c echo.Context) error {
		return c.NoContent(fiber.StatusNoContent)
	})

	return echoServer
}
