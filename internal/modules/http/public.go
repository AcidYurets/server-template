package http

import (
	"server-template/internal/pkg/routers"
)

func publicRouter(router routers.Router) PublicRouter {
	router.Static("/docs", "./docs")

	return router
}
