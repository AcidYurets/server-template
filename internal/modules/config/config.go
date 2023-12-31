package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"time"
)

// Config ...
type Config struct {
	DBUser    string `envconfig:"DB_USER"`
	DBPass    string `envconfig:"DB_PASS"`
	DBHost    string `envconfig:"DB_HOST" default:"db"`
	DBPort    string `envconfig:"DB_PORT" default:"5432"`
	DBName    string `envconfig:"DB_NAME" default:"server_db"`
	DBSSLMode string `envconfig:"DB_SSL_MODE" default:"disable" validate:"oneof=disable enable"`

	SQLSlowThreshold int    `envconfig:"SQL_SLOW_THRESHOLD" default:"600"`
	TraceSQLCommands bool   `envconfig:"TRACE_SQL_COMMANDS" default:"false"`
	AutoMigrate      bool   `envconfig:"AUTO_MIGRATE" default:"false"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info" validate:"oneof=debug info warn error dpanic panic fatal"`

	HTTPServerHost         string        `envconfig:"HTTP_SERVER_HOST" default:"0.0.0.0"`
	HTTPServerPort         string        `envconfig:"HTTP_SERVER_PORT" default:"8080"`
	HTTPServerReadTimeOut  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"10m"`
	HTTPServerWriteTimeOut time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"13m"`

	DevMode bool `envconfig:"DEV_MODE" default:"false"`
}

func NewConfig(logger *zap.Logger, logLevel zap.AtomicLevel) (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, err
	}

	logger.Info("Получена конфигурация", zap.Any("config", config))

	// Принудительно инициализируем уровень логирования из конфигурации
	err = logLevel.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		return Config{}, err
	}

	return config, err
}
