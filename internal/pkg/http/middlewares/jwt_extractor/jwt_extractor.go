package jwt_extractor

import (
	"fmt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// New ...
func New(config ...Config) echo.MiddlewareFunc {
	cfg := makeCfg(config)

	extractors, err := middleware.CreateExtractors(cfg.TokenLookup)
	if err != nil {
		panic(fmt.Errorf("не удалось создать массив extractors: %w", err))
	}

	// Return middleware handler
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Filter request to skip middleware
			if cfg.Filter != nil && cfg.Filter(c) {
				return next(c)
			}

			var lastExtractorErr error

			for _, extractor := range extractors {
				auths, extrErr := extractor(c)
				if extrErr != nil {
					lastExtractorErr = extrErr
					continue
				}
				for _, auth := range auths {
					token := auth
					// Store user information from token into context.
					c.Set(cfg.ContextKey, token)

					err := cfg.SuccessHandler(c)
					if err != nil {
						return err
					}

					return next(c)
				}
			}

			var err error
			if lastExtractorErr != nil {
				err = &echojwt.TokenExtractionError{Err: lastExtractorErr}
			}

			tmpErr := cfg.ErrorHandler(c, err)
			return tmpErr
		}
	}
}
