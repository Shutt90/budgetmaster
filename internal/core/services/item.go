package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type ItemService struct {
	itemRepository ports.ItemRepository
	clock          ports.Clock
}

func NewItemService(ir ports.ItemRepository, clock ports.Clock) *ItemService {
	if ir == nil || clock == nil {
		fmt.Printf("itemRepo:%#v\nclock:%#v\n", ir, clock)
		panic("unable to start new item repository")
	}

	return &ItemService{
		itemRepository: ir,
		clock:          clock,
	}
}

func (is *ItemService) Create(i domain.Item) error {
	if err := is.itemRepository.Create(i); err != nil {
		return err
	}

	log.Infof("item bought: %s\ncost: %d", i.Name, i.Cost)

	return nil
}

func (is *ItemService) GetDefaultMonthlyItems() ([]domain.Item, error) {
	now := is.clock.Now()
	month := now.Month()
	year := now.Year()

	items, err := is.itemRepository.GetMonthlyItems(int(month), year)
	if err != nil {
		return []domain.Item{}, err
	}

	var itemsToShow []domain.Item
	for _, item := range items {
		item.CreatedAt = nil
		item.UpdatedAt = nil
		item.RemovedOccuringAt = nil
		item.RemovedOccuringAt = nil
		item.Cost = 0
		item.ID = 0

		itemsToShow = append(itemsToShow, item)
	}

	log.Info(itemsToShow)

	return itemsToShow, nil
}

func (is *ItemService) GetMonthlyItems(month string, year string) ([]domain.Item, error) {
	y, err := strconv.Atoi(year)
	if err != nil {
		log.Error(err)
		return []domain.Item{}, err
	}
	log.Info(string(month[0]) + month[1:])
	m, err := time.Parse(time.January.String(), strings.ToUpper(string(month[0]))+month[1:])
	if err != nil {
		return []domain.Item{}, echo.ErrBadRequest
	}

	items, err := is.itemRepository.GetMonthlyItems(int(m.Month()), y)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (is *ItemService) SwitchRecurringPayments(id string, isRecurring bool) error {
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	if err := is.itemRepository.SwitchRecurringPayments(idUint, isRecurring); err != nil {
		return err
	}

	return nil
}
