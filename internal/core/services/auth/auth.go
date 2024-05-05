package auth

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	template "github.com/Shutt90/budgetmaster/templating"
)

func NewWithConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			f := template.NewFlash("unauthorized", true)
			return c.Render(http.StatusUnauthorized, "flash", f)
		},
	})

}
