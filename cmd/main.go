package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
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
		err := c.Render(200, "index", "")
		log.Error(err)
		return c.Render(200, "index", "")
	})

	r.Router.POST("/login", func(c echo.Context) error {
		if err := h.Login(c); err != nil {
			log.Error(err)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			f := template.NewFlash("unable to login", true)
			c.Render(http.StatusForbidden, "flash", f)

			return err
		}

		return c.Render(http.StatusAccepted, "logged", "success")
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

	g := e.Group("/item")
	u := e.Group("/user")

	g.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			log.Error(err)
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		},
	}))

	u.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			log.Error(err)
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		},
	}))

	g.GET("/monthly", h.GetMonth)
	g.POST("/create", h.CreateItem)
	g.PATCH("/:id", h.SwitchRecurring)

	u.PATCH("/user/:id", h.ChangePassword)

	e.Logger.Fatal(e.Start(":9002"))
}
