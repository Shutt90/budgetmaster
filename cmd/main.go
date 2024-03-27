package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/Shutt90/budgetmaster/internal/repositories"
)

func main() {
	e := echo.New()
	e.Static("/public/css", "css")
	e.Static("/public/images", "images")

	db := repositories.NewDB(os.Getenv("DSN"))

	clock := services.NewClock()
	itemService := repositories.NewItemRepository(db, clock)

	itemService.Get(1)

	if err := itemService.CreateItemTable(); err != nil {
		fmt.Println(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Budget Master")
	})

	e.Logger.Fatal(e.Start(":9001"))
}
