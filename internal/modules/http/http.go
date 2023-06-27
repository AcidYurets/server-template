package http

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"server-template/internal/pkg/routers"
)

type Routers struct {
	fx.Out

	Api    ApiRouter
	Public PublicRouter
}

type ApiRouter routers.Router
type PublicRouter routers.Router

// NewHTTP создает http все необходимые роутеры
func NewHTTP(server *echo.Echo) (Routers, error) {
	return Routers{
		Public: publicRouter(&routers.EchoRouter{Router: server}),
		Api:    apiRouter(&routers.EchoRouter{Router: server.Group("/api")}),
	}, nil
}
