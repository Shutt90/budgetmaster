package repositories

import (
	"database/sql"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type mockUserRepository struct {
	*sql.DB
}

func (ur *mockUserRepository) GetByEmail(email string) (*domain.User, error) {
	return ur.GetByEmail(email)
}

func (ur *mockUserRepository) ChangePassword(id uint64, email string, password string) error {
	return ur.ChangePassword(id, email, password)
}
