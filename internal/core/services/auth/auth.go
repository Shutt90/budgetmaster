package auth

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/Shutt90/budgetmaster/internal/core/ports"
	template "github.com/Shutt90/budgetmaster/templating"
)

type JwtService struct {
	jwt ports.JwtIface
}

func New(jwt ports.JwtIface) *JwtService {
	return &JwtService{
		jwt: jwt,
	}

}

func (j *JwtService) SetConfig() echo.MiddlewareFunc {
	return j.jwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			f := template.NewFlash("unauthorized", true)
			return c.Render(http.StatusUnauthorized, "flash", f)
		},
	})
}
