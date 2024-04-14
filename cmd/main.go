package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	g := e.Group("admin")

	r.Router.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", "")
	})

	r.Router.POST("/login", func(c echo.Context) error {
		if err := h.Login(c); err != nil {
			g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
				return false, nil
			}))
		}

		g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			return true, nil
		}))

		return c.Render(http.StatusAccepted, "logged", "success")
	})

	g.GET("/items", func(c echo.Context) error {
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

	g.GET("/item/monthly", h.GetMonth)
	g.POST("/item/create", h.CreateItem)
	g.PATCH("/item/:id", h.SwitchRecurring)

	g.PATCH("/login/user/:id", h.ChangePassword)

	e.Logger.Fatal(e.Start(":9002"))
}
