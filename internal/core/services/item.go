package services

import (
	"fmt"

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

func (is *ItemService) GetMonthlyItems(month string, year int) ([]domain.Item, error) {
	items, err := is.itemRepository.GetMonthlyItems(month, year)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (is *ItemService) SwitchRecurringPayments(id uint64, isRecurring bool) error {
	if err := is.itemRepository.SwitchRecurringPayments(id, isRecurring); err != nil {
		return err
	}

	return nil
}
