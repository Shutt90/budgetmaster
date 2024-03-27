package services

import (
	"time"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type itemService struct {
	itemRepository ports.ItemRepository
}

func New(ir ports.ItemRepository) *itemService {
	return &itemService{
		itemRepository: ir,
	}
}

func (srv *itemService) Create(i domain.Item) error {
	if err := srv.itemRepository.Create(i); err != nil {
		return err
	}

	return nil
}

func (srv *itemService) GetDefaultMonthlyItems() ([]domain.Item, error) {
	now := time.Now()
	month := now.Month()
	year := now.Year()

	items, err := srv.itemRepository.GetMonthlyItems(month.String(), year)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (srv *itemService) GetMonthlyItems(month string, year int) ([]domain.Item, error) {
	items, err := srv.itemRepository.GetMonthlyItems(month, year)
	if err != nil {
		return []domain.Item{}, err
	}

	return items, nil
}

func (srv *itemService) SwitchRecurringPayments(id uint64, isRecurring bool) error {
	if err := srv.itemRepository.SwitchRecurringPayments(id, isRecurring); err != nil {
		return err
	}

	return nil
}
