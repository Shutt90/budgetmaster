package repositories

import (
	"database/sql"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
)

type userRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (ur *userRepository) GetByEmail(email string) (*domain.User, error) {
	row := ur.DB.QueryRow("SELECT id, firstName, surname, password FROM user WHERE email = ?;", email)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, row.Err()
	}

	u := domain.User{}
	row.Scan(
		&u.ID,
		&u.FirstName,
		&u.Surname,
		&u.Password,
	)

	u.Email = email

	return &u, nil
}

func (ur *userRepository) ChangePassword(email string, password string) error {
	_, err := ur.DB.Exec("UPDATE user SET password TO ? WHERE email = ?;", password, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
	}

	return nil
}
