package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("checks user created", func(t *testing.T) {
		u := NewUser(
			"firstName",
			"surname",
			"email",
			"password",
			[]string{"admin", "user"},
		)

		if u == nil {
			t.Error("new user not created")
		}

		if u.FirstName != "firstName" {
			t.Error("first name incorrect")
		}

		if u.Surname != "surname" {
			t.Error("surname incorrect")
		}

		if u.Email != "email" {
			t.Error("email incorrect")
		}

		if u.Password != "password" {
			t.Error("password incorrect")
		}

		if u.Roles == nil {
			t.Error("roles incorrect")
		}
	})
}
