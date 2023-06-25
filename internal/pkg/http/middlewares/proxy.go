package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"server-template/internal/pkg/constants"
	"server-template/internal/pkg/http/error_handler/http_errors"
	"server-template/internal/utils"
)

type Proxy struct {
	Host      string
	Port      string
	SecretKey string
	Rewrite   map[string]string
}

func ProxyHandler(cfg Proxy) echo.MiddlewareFunc {
	url, err := utils.NewURL(cfg.Host, cfg.Port)
	if err != nil {
		// Паникуем, если не получилось создать url по данным из конфигурации
		panic(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set(constants.HeaderAuthorization, cfg.SecretKey)

			proxyMiddleware := middleware.ProxyWithConfig(middleware.ProxyConfig{
				Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: url}}),
				Rewrite:  cfg.Rewrite,
				ModifyResponse: func(resp *http.Response) error {
					if resp.StatusCode >= 400 {
						decoder := json.NewDecoder(resp.Body)
						var errResp http_errors.ErrorResponse
						if err := decoder.Decode(&errResp); err != nil {
							return fmt.Errorf("ошибка обработки ответа от проксируемого сервера: %w", err)
						}
						errResp.Code = resp.StatusCode

						return &errResp
					}

					return nil
				},
			})

			return proxyMiddleware(next)(c)
		}
	}
}
