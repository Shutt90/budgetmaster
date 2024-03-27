package ports

import (
	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type ItemRepository interface {
	CreateItemTable() error
	Create(domain.Item) error
	Get(uint64) (domain.Item, error)
	GetMonthlyItems(string, int) ([]domain.Item, error)
	SwitchRecurringPayments(uint64, bool) error
}
