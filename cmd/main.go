package main

import (
	"net/http"
	"os"

	"github.com/Shutt90/budgetmaster/internal/repositories"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/public/css", "css")
	e.Static("/public/images", "images")

	db := repositories.New("libsql://budgetmaster.turso.io", os.Getenv("DB_PASS"))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Budget Master")
	})

	e.Logger.Fatal(e.Start(":9001"))
}
