package services

import (
	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type itemService struct {
	itemRepository ports.ItemRepository
	clock          ports.Clock
}

func NewItemService(ir ports.ItemRepository) *itemService {
	return &itemService{
		itemRepository: ir,
	}
}

func (is *itemService) Create(i domain.Item) error {
	if err := is.itemRepository.Create(i); err != nil {
		return err
	}

	return nil
}

func (is *itemService) GetDefaultMonthlyItems() ([]domain.Item, error) {
	now := is.clock.Now()
	month := now.Month()
	year := now.Year()

	items, err := is.itemRepository.GetMonthlyItems(month.String(), year)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (is *itemService) GetMonthlyItems(month string, year int) ([]domain.Item, error) {
	items, err := is.itemRepository.GetMonthlyItems(month, year)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (is *itemService) SwitchRecurringPayments(id uint64, isRecurring bool) error {
	if err := is.itemRepository.SwitchRecurringPayments(id, isRecurring); err != nil {
		return err
	}

	return nil
}
