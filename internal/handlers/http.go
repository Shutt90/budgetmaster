package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ErrNotProcessable = errors.New("unprocessable entity")
	ErrBadRequest     = errors.New("bad request")
	ErrInternal       = errors.New("internal server error")
)

type HTTPHandler struct {
	is  *services.ItemService
	us  *services.UserService
	jwt *jwtCustomClaims
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func NewHttpHandler(is *services.ItemService, us *services.UserService) *HTTPHandler {
	return &HTTPHandler{
		is: is,
		us: us,
	}
}

func (h *HTTPHandler) ItemService() *services.ItemService {
	return h.is
}

func (h *HTTPHandler) UserService() *services.UserService {
	return h.us
}

func (h *HTTPHandler) GetDefaults(c echo.Context) ([]domain.Item, error) {
	items, err := h.is.GetDefaultMonthlyItems()
	if err != nil {
		c.JSON(500, ErrBadRequest)

		return []domain.Item{}, err
	}

	return items, nil
}

func (h *HTTPHandler) GetMonth(c echo.Context) error {
	month, year := c.QueryParam("month"), c.QueryParam("year")
	items, err := h.is.GetMonthlyItems(month, year)
	if err != nil {
		c.JSON(500, ErrBadRequest)

		return err
	}

	c.JSON(200, items)

	return nil
}

func (h *HTTPHandler) CreateItem(c echo.Context) error {
	cost, err := strconv.ParseUint(strings.ReplaceAll(c.FormValue("cost"), ".", ""), 10, 64)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, ErrBadRequest)
	}
	year, err := strconv.ParseUint(c.FormValue("year"), 10, 16)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, ErrBadRequest)
	}
	isRecurring, err := strconv.ParseBool(c.FormValue("isRecurring"))
	if err != nil {
		isRecurring = false
	}

	err = h.is.Create(
		domain.NewItem(
			c.FormValue("name"),
			c.FormValue("description"),
			c.FormValue("location"),
			c.FormValue("month"),
			uint16(year),
			cost,
			isRecurring,
		),
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrNotProcessable)

		return err
	}

	c.JSON(200, "success")

	return nil
}

func (h *HTTPHandler) SwitchRecurring(c echo.Context) error {
	id := c.Param("id")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		log.Error(err)

		return err
	}

	var isRecurring bool
	err = json.Unmarshal(body, &isRecurring)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrNotProcessable)
		log.Error(err)

		return err
	}

	h.is.SwitchRecurringPayments(id, isRecurring)

	return nil
}

func (h *HTTPHandler) Login(c echo.Context) error {
	_, err := h.us.Login(c.FormValue("email"), c.FormValue("password"))
	if err != nil {
		c.JSON(http.StatusNotFound, echo.ErrNotFound)
		log.Error(err)

		return err
	}

	c.JSON(http.StatusAccepted, "accepted")

	h.jwt.Name = c.FormValue("email")
	h.jwt.Admin = true
	h.jwt.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, h.jwt)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (h *HTTPHandler) ChangePassword(c echo.Context) error {
	id := c.Param("id")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, echo.ErrUnprocessableEntity)
		log.Error(err)

		return err
	}

	u := domain.User{}
	if err := json.Unmarshal(body, &u); err != nil {
		c.JSON(http.StatusBadRequest, echo.ErrBadRequest)
		log.Error(err)

		return err
	}

	if err := h.us.ChangePassword(id, u.Email, u.Password); err != nil {
		c.JSON(http.StatusNotFound, echo.ErrNotFound)
		log.Error(err)

		return err
	}

	c.JSON(http.StatusAccepted, "accepted")

	return nil
}
