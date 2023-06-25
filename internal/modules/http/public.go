package http

import (
	"server-template/internal/pkg/routers"
)

func publicRouter(echoRouter routers.IEchoRouter) PublicRouter {
	// Тут middleware для /

	return &routers.EchoRouter{
		Router: echoRouter,
	}
}
