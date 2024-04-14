package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/Shutt90/budgetmaster/internal/handlers"
	"github.com/Shutt90/budgetmaster/internal/repositories"
	"github.com/Shutt90/budgetmaster/internal/router"
	template "github.com/Shutt90/budgetmaster/templating"
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

	if err := ur.CreateUserTable(); err != nil {
		log.Error("tried to create db but couldnt: ", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Static("/public/css", "css")
	e.Static("/public/images", "images")
	e.Renderer = template.NewTemplate()

	h := handlers.NewHttpHandler(itemService, userService)
	r := router.New(e)

	r.Router.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", "")
	})

	r.Router.GET("/items", func(c echo.Context) error {
		items, err := h.GetDefaults(c)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		err = c.Render(200, "items", items)
		if err != nil {
			log.Error(err)
		}
		return c.Render(200, "items", items)
	})
	r.Router.GET("/item/monthly", h.GetMonth)
	r.Router.POST("/item/create", h.CreateItem)
	r.Router.PATCH("/item/:id", h.SwitchRecurring)

	r.Router.POST("/login", h.Login)
	r.Router.PATCH("/login/user/:id", h.ChangePassword)

	e.Logger.Fatal(e.Start(":9002"))
}
