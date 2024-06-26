package services

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/repositories"
)

func TestLogin(t *testing.T) {
	type testcase struct {
		name        string
		expected    domain.User
		expectedErr error
	}

	testcases := []testcase{
		{
			name: "check login success",
			expected: domain.User{
				ID:        1,
				Email:     "test@example.com",
				Password:  "password",
				FirstName: "fname",
				Surname:   "surname",
				Roles:     []string{"admin", "user"},
			},
			expectedErr: nil,
		},
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockService := NewUserService(repositories.NewUserRepository(db), NewMockCrypt())

	userMockRows := sqlmock.NewRows([]string{"id", "firstName", "surname", "password", "roles"}).
		AddRow("1", "fname", "surname", "password", "{admin,user}")

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, firstName, surname, password, roles FROM user WHERE email = ?;`)).
		WithArgs(
			"test@example.com",
		).WillReturnRows(userMockRows)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := mockService.Login(tc.expected.Email, tc.expected.Password)
			if err != tc.expectedErr {
				t.Errorf("unexpected error in test %s\nwant: %s\ngot: %s", tc.name, tc.expectedErr.Error(), err.Error())
			}

			diff := cmp.Diff(u, tc.expected)
			if diff != "" {
				t.Errorf("unexpected user in test %s\n%s", tc.name, diff)
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
				t.Errorf("unexpected error in test %s\nwant: %s\ngot: %s", tc.name, tc.expectedErr.Error(), err.Error())
			}
		})
	}
}
