package ports

import (
	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type ItemRepository interface {
	CreateItemTable()
	Create(domain.Item) error
	Get(uint64) *domain.Item
	GetMonthlyItems(string, int) ([]domain.Item, error)
	SwitchRecurringPayments(uint64, bool) error
}
