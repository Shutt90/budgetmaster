package repositories

import (
	"database/sql"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type mockRepository struct {
	mock  *sql.DB
	clock ports.Clock
}

func (mr *mockRepository) CreateItemTable() error {
	return mr.CreateItemTable()
}

func (mr *mockRepository) Create(i domain.Item) error {
	return mr.Create(i)
}

func (mr *mockRepository) Get(id uint64) (domain.Item, error) {
	return mr.Get(id)
}

func (mr *mockRepository) GetMonthlyItems(m string, y int) error {
	return mr.GetMonthlyItems(m, y)
}

func (mr *mockRepository) SwitchRecurringPayments(id uint64, isRecurring bool) error {
	return mr.SwitchRecurringPayments(id, isRecurring)
}
