package repositories

import (
	"database/sql"
	"os"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/labstack/gommon/log"
)

type userRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (db *userRepository) CreateUserTable() error {
	db.DB.Begin()

	queryBytes, err := os.ReadFile("internal/migrations/user_table_schema.sql")
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = db.Exec(string(queryBytes))
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (ur *userRepository) GetByEmail(email string) (domain.User, error) {
	row := ur.DB.QueryRow("SELECT id, firstName, surname, password FROM user WHERE email = ?;", email)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			log.Error(row.Err())
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, row.Err()
	}

	u := domain.User{}
	row.Scan(
		&u.ID,
		&u.FirstName,
		&u.Surname,
		&u.Password,
		&u.IsAdmin,
	)

	u.Email = email

	return u, nil
}

func (ur *userRepository) ChangePassword(id uint64, email string, password string) error {
	_, err := ur.DB.Exec("UPDATE user SET password TO ? WHERE email = ? AND id = ?;", password, email, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
	}

	return nil
}
