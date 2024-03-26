package ports

import (
	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type Database interface {
	Get(id string) *domain.Item
}
