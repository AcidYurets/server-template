package http

import (
	"server-template/internal/pkg/routers"
)

func apiRouter(router routers.Router) ApiRouter {
	return router
}
