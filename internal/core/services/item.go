package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
	"github.com/labstack/gommon/log"
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

	return nil
}

func (is *ItemService) GetDefaultMonthlyItems() ([]domain.Item, error) {
	now := is.clock.Now()
	month := now.Month()
	year := now.Year()

	items, err := is.itemRepository.GetMonthlyItems(month.String(), year)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (is *ItemService) GetMonthlyItems(month string, year string) ([]domain.Item, error) {
	y, err := strconv.Atoi(year)
	if err != nil {
		log.Error(err)
		return []domain.Item{}, err
	}

	m := strings.ToUpper(string(month[0])) + month[1:]

	items, err := is.itemRepository.GetMonthlyItems(m, y)
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
