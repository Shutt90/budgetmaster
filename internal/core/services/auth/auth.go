package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	template "github.com/Shutt90/budgetmaster/templating"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
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

func (claims *JwtCustomClaims) SetLoggedIn(name string) {
	*claims = JwtCustomClaims{
		Name:  name,
		Admin: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
}

func (claims *JwtCustomClaims) SetAdmin(name string) {
	*claims = JwtCustomClaims{
		Name:  name,
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
}

func (claims *JwtCustomClaims) GetClaims(c echo.Context, t TokenString) (jwt.Claims, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}

	jwtVal := token.Value
	if jwtVal == "" {
		return nil, nil
	}

	parsedClaims, err := jwt.ParseWithClaims(jwtVal, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedClaims.Claims.(*JwtCustomClaims); ok && parsedClaims.Valid {
		return claims, nil
	} else {
		log.Error("custom claims not valid")
		return nil, errors.New("custom claims not valid")
	}
}
