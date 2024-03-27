package ports

import (
	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type ItemRepository interface {
	Get(id string) *domain.Item
}
