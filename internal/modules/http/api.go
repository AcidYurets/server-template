package http

import (
	"server-template/internal/pkg/routers"
)

func apiRouter(echoRouter routers.IEchoRouter) ApiRouter {
	// Тут middleware для /api

	return &routers.EchoRouter{
		Router: echoRouter,
	}
}
