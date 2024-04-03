package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/Shutt90/budgetmaster/internal/handlers"
	"github.com/Shutt90/budgetmaster/internal/repositories"
)

func main() {
	db := repositories.NewDB(os.Getenv("DATABASE"), os.Getenv("TOKEN"))

	clock := services.NewClock()
	crypt := services.NewCrypt()
	ir := repositories.NewItemRepository(db, clock)
	ur := repositories.NewUserRepository(db)
	itemService := services.NewItemService(ir, clock)
	userService := services.NewUserService(ur, crypt)

	if err := ir.CreateItemTable(); err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	e.Static("/public/css", "css")
	e.Static("/public/images", "images")

	h := handlers.NewHttpHandler(itemService, userService)

	e.GET("/", func(c echo.Context) error {
		return h.GetDefaults(c)
	})

	e.Logger.Fatal(e.Start(":9002"))
}
