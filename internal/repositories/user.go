package repositories

import (
	"database/sql"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type userRepository struct {
	*sql.DB
	clock ports.Clock
}

func (ur *userRepository) GetByLogin(email string, password []byte) (*domain.User, error) {
	row := ur.DB.QueryRow("SELECT id, firstName, lastName FROM user WHERE email = ? AND password = ?", email, string(password))
	if row.Err() == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	u := domain.User{}
	row.Scan(
		&u.ID,
		&u.FirstName,
		&u.Surname,
	)

	u.Email = email

	return &u, nil
}

func (ur *userRepository) ChangePassword(email string, password []byte) error {
	_, err := ur.DB.Exec("UPDATE password FROM user WHERE email = ? AND password = ?", email, string(password))
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
	}

	return nil
}
