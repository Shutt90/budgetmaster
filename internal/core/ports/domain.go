package ports

import (
	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type ItemRepository interface {
	CreateItemTable() error
	Create(domain.Item) error
	Get(uint64) (domain.Item, error)
	GetMonthlyItems(int, int) ([]domain.Item, error)
	SwitchRecurringPayments(uint64, bool) error
}

type UserRepository interface {
	GetByEmail(string) (domain.User, error)
	ChangePassword(uint64, string, string) error
}
