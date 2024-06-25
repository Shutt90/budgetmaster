package domain

import (
	"database/sql"
)

type User struct {
	ID        uint64
	FirstName string
	Surname   string
	Email     string
	Password  string
	Roles     []string
	CreatedAt *sql.NullTime `json:"createdAt,omitempty"`
	UpdatedAt *sql.NullTime `json:"updatedAt,omitempty"`
}

func NewUser(fname, surname, email, password string, roles []string) *User {
	return &User{
		FirstName: fname,
		Surname:   surname,
		Email:     email,
		Password:  password,
		Roles:     roles,
	}
}
