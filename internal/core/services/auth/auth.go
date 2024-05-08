package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	template "github.com/Shutt90/budgetmaster/templating"
)

func SetConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			f := template.NewFlash("unauthorized", true)
			return c.Render(http.StatusUnauthorized, "flash", f)
		},
	})
}

func GetClaims(c echo.Context) string {
	_, err := c.Cookie("token")
	if err != nil {
		log.Error(err)
		return ""
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)["username"].(string)

	log.Info("claims", claims)

	return claims
}
