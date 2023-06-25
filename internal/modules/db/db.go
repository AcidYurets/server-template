package db

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"server-template/internal/modules/config"
	"server-template/internal/modules/db/ent"

	_ "server-template/internal/modules/db/ent/runtime"
)

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --target ./ent --feature  sql/execquery,sql/upsert,intercept,schema/snapshot ./schema

func NewDBClient(cfg config.Config, logger *zap.Logger) (*ent.Client, error) {
	client, err := connectDB(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	return client, nil
}

func InvokeDBClient(
	client *ent.Client,
	cfg config.Config,
	lifecycle fx.Lifecycle,
) error {
	ctx := context.Background()

	// Делаем миграцию при необходимости
	if cfg.AutoMigrate {
		if err := client.Schema.Create(ctx); err != nil {
			return fmt.Errorf("ошибка при миграции схемы: %w", err)
		}
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return nil
}
