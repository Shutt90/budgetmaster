package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ErrNotProcessable = errors.New("unprocessable entity")
	ErrBadRequest     = errors.New("bad request")
	ErrInternal       = errors.New("internal server error")
)

type HTTPHandler struct {
	is *services.ItemService
	us *services.UserService
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

func (h *HTTPHandler) GetDefaults(c echo.Context) error {
	items, err := h.is.GetDefaultMonthlyItems()
	if err != nil {
		c.JSON(500, ErrBadRequest)

		return err
	}

	c.JSON(200, items)

	return nil
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
	err := h.is.Create(c.Request().Body)
	if err != nil {
		c.JSON(500, ErrBadRequest)

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

	if err := h.us.Login(u.Email, u.Password); err != nil {
		c.JSON(http.StatusNotFound, echo.ErrNotFound)
		log.Error(err)

		return err
	}

	c.JSON(http.StatusAccepted, "accepted")

	return nil
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
