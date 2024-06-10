package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	template "github.com/Shutt90/budgetmaster/templating"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func NewJwtCustomClaims(name string, admin bool) *JwtCustomClaims {
	return &JwtCustomClaims{}
}

type TokenString string

func NewTokenString(t string) TokenString {
	return TokenString(t)
}

func SetConfig(t TokenString) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(t),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			f := template.NewFlash("unauthorized", true)
			return c.Render(http.StatusUnauthorized, "flash", f)
		},
	})
}

func (t TokenString) SetLoggedIn(name string) {
	NewJwtCustomClaims(name, false)
}

func (t TokenString) SetAdmin(name string) {
	NewJwtCustomClaims(name, true)
}

func (t TokenString) GetClaims(c echo.Context) string {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(string(t), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t), nil
	})
	log.Info(c.Cookie("token"))

	if err := claims.Valid(); err != nil {
		log.Error("unauthorised attempt to login")
		return ""
	}

	log.Info(claims["name"])
	return ""
}

func (t TokenString) SetClaims(c echo.Context) string {
	var claims JwtCustomClaims
	jwt.ParseWithClaims(string(t), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t), nil
	})
	log.Info(c.Cookie("token"))

	if err := claims.Valid(); err != nil {
		log.Error("unauthorised attempt to login")
		return ""
	}

	log.Info(claims.Name])
	return ""
}
