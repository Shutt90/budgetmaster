package ports

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtIface interface {
	WithConfig(config echojwt.Config) echo.MiddlewareFunc
}
