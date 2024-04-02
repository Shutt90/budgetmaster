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

func (ur *mockUserRepository) ChangePassword(email string, password string) error {
	return ur.ChangePassword(email, password)
}
