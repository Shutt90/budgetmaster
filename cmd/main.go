package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/Shutt90/budgetmaster/internal/handlers"
	"github.com/Shutt90/budgetmaster/internal/repositories"
	"github.com/Shutt90/budgetmaster/internal/router"
)

func main() {
	godotenv.Load()
	db := handlers.NewDB(os.Getenv("DATABASE"), os.Getenv("TOKEN"))

	clock := services.NewClock()
	crypt := services.NewCrypt()
	ir := repositories.NewItemRepository(db, clock)
	ur := repositories.NewUserRepository(db)
	itemService := services.NewItemService(ir, clock)
	userService := services.NewUserService(ur, crypt)

	if err := ir.CreateItemTable(); err != nil {
		log.Error("tried to create db but couldnt: ", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Static("/public/css", "css")
	e.Static("/public/images", "images")

	h := handlers.NewHttpHandler(itemService, userService)
	r := router.New(e)

	r.Router.GET("/items", h.GetDefaults)
	r.Router.GET("/item/monthly", h.GetMonth)
	r.Router.POST("/item/create", h.CreateItem)
	r.Router.PATCH("/item/:id", h.SwitchRecurring)

	e.Logger.Fatal(e.Start(":9002"))
}
