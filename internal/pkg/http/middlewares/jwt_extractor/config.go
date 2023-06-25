package jwt_extractor

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	defaultTokenLookup = "header:Authorization:Bearer "
)

var (
	ErrMissingJWT = errors.New("токен отсутствует")
)

// Config defines the config for JWT middleware
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(c echo.Context) bool

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler func(c echo.Context) error

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: 401 Invalid or expired JWT
	ErrorHandler func(c echo.Context, err error) error

	// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// 			`<cut-prefix>` is argument value to cut/trim prefix of the extracted value. This is useful if header
	//			value has static prefix like `Authorization: <auth-scheme> <authorisation-parameters>` where part that we
	//			want to cut is `<auth-scheme> ` note the space at the end.
	//			In case of JWT tokens `Authorization: Bearer <token>` prefix we cut is `Bearer `.
	// If prefix is left empty the whole value is returned.
	// - "query:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	// - "form:<name>"
	// Multiple sources example:
	// - "header:Authorization:Bearer ,cookie:myowncookie"
	TokenLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default: "Bearer".
	AuthScheme string

	// Context key to store user information from the token into context.
	// Optional. Default: "token".
	ContextKey string
}

// makeCfg function will check correctness of supplied configuration
// and will complement it with default values instead of missing ones
func makeCfg(config []Config) (cfg Config) {
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c echo.Context) error { return nil }
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c echo.Context, err error) error {
			// TODO Подумать над обработкой
			if err.Error() == "Missing or malformed JWT" {
				c.Response().WriteHeader(http.StatusBadRequest)
				return ErrMissingJWT
			}
			return c.String(http.StatusUnauthorized, "Invalid or expired JWT")
		}
	}

	if cfg.TokenLookup == "" {
		cfg.TokenLookup = defaultTokenLookup
	}
	if cfg.AuthScheme == "" {
		cfg.AuthScheme = "Bearer"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "token"
	}
	return cfg
}
