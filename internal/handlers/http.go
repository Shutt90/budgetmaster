package handlers

import (
	"errors"

	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/labstack/echo/v4"
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
