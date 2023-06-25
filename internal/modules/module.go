package modules

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"server-template/internal/modules/config"
	"server-template/internal/modules/db"
	"server-template/internal/modules/domain"
	"server-template/internal/modules/echo"
	"server-template/internal/modules/http"
	"server-template/internal/modules/logger"
)

var (
	AppModule = fx.Options(
		logger.Module,
		config.Module,
		db.Module,
		echo.Module,
		http.Module,

		domain.Module,

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)

	AppInvokables = fx.Options(
		logger.Invokables,
		config.Invokables,
		db.Invokables,

		domain.Invokables,
	)
)
