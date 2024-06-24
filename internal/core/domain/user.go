package domain

import (
	"database/sql"
)

type User struct {
	ID        uint64
	FirstName string
	Surname   string
	Email     string
	LoggedIn  bool
	Password  string
	Roles     []string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func NewUser(fname, surname, email, password string) *User {
	return &User{
		FirstName: fname,
		Surname:   surname,
		Email:     email,
		Password:  password,
	}
}
