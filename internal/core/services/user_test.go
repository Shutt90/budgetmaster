package services

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/Shutt90/budgetmaster/internal/repositories"
)

func TestLogin(t *testing.T) {
	type testcase struct {
		name           string
		email          string
		password       string
		expectedErr    error
		expectedResult string
	}

	testcases := []testcase{
		{
			name:        "check login success",
			email:       "test@example.com",
			password:    "password",
			expectedErr: nil,
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockService := NewUserService(repositories.NewUserRepository(db), NewMockCrypt())

	userMockRows := sqlmock.NewRows([]string{"id", "firstName", "surname", "password"}).
		AddRow("1", "fname", "surname", "password")

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, firstName, surname, password FROM user WHERE email = ?;`)).
		WithArgs(
			"test@example.com",
		).WillReturnRows(userMockRows)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := mockService.Login(tc.email, tc.password)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\nwant: %s\ngot: %s", tc.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type testcase struct {
		name           string
		email          string
		password       string
		id             string
		expectedErr    error
		expectedResult string
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockService := NewUserService(repositories.NewUserRepository(db), NewMockCrypt())

	testcases := []testcase{
		{
			name:        "get user success",
			email:       "test@example.com",
			password:    "password",
			id:          "1",
			expectedErr: nil,
		},
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE user SET password TO ? WHERE email = ?;`)).
		WithArgs(
			"test@example.com",
			"password",
		).WillReturnResult(sqlmock.NewResult(1, 1))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := mockService.ChangePassword(tc.id, tc.email, tc.password)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\nwant: %s\ngot: %s", tc.expectedErr.Error(), err.Error())
			}
		})
	}
}
