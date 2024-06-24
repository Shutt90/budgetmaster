package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/Shutt90/budgetmaster/internal/core/services/auth"
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
	t := auth.NewTokenString(os.Getenv("JWT_SECRET"))
	authConfig := auth.SetConfig(t)
	jwt := &auth.JwtCustomClaims{}

	var migrateFlag bool
	flag.BoolVar(&migrateFlag, "migrate", false, "should a database be migrated from scratch")
	flag.Parse()

	if migrateFlag {
		if err := ir.CreateItemTable(); err != nil {
			log.Error("tried to create db but couldnt: ", err)
			panic(err)
		}

		if err := ur.CreateUserTable(); err != nil {
			log.Error("tried to create db but couldnt: ", err)
			panic(err)
		}

		if err := ur.AddDebugUser(); err != nil {
			log.Error("tried to add debug user but couldnt: ", err)
			panic(err)
		}
	}

	e := echo.New()
	e.Static("/public/css", "css")
	e.Static("/public/images", "images")
	e.Renderer = template.NewTemplate()
	h := handlers.NewHttpHandler(itemService, userService, jwt)
	r := router.New(e)

	r.Router.GET("/", func(c echo.Context) error {
		_, err := h.JwtService().GetClaims(c, t)
		if err != nil {
			return c.Render(http.StatusUnauthorized, "index", nil)
		}

		return c.Render(http.StatusOK, "index", h.JwtService())
	})

	r.Router.POST("/login", func(c echo.Context) error {
		if err := h.Login(c); err != nil {
			f := template.NewFlash("unable to login", true)
			return c.Render(http.StatusForbidden, "flash", f)
		}

		return c.Render(http.StatusAccepted, "index", "success")
	})

	g := e.Group("/item")
	u := e.Group("/user")

	g.Use(authConfig)
	u.Use(authConfig)

	r.Router.GET("/items", func(c echo.Context) error {
		items, err := h.GetDefaults(c)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return err
		}

		return c.Render(200, "items", items)
	})

	g.GET("/monthly", h.GetMonth)
	g.POST("/create", h.CreateItem)
	g.PATCH("/:id", h.SwitchRecurring)

	u.PATCH("/user/:id", h.ChangePassword)

	e.Logger.Fatal(e.Start(":9002"))
}
