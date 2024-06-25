package repositories

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

func TestGetByEmail(t *testing.T) {
	type testcase struct {
		name           string
		email          string
		ur             ports.UserRepository
		expectedErr    error
		expectedResult string
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockRepo := NewUserRepository(db)

	testcases := []testcase{
		{
			name:           "get user success",
			email:          "test@example.com",
			ur:             mockRepo,
			expectedErr:    nil,
			expectedResult: `{"ID":1,"FirstName":"fname","Surname":"surname","Email":"test@example.com","Password":"password","Roles":["admin","user"]}`,
		},
		{
			name:           "get user not found",
			email:          "test2@example.com",
			ur:             mockRepo,
			expectedErr:    ErrNotFound,
			expectedResult: "null",
		},
	}

	userMockRows := sqlmock.NewRows([]string{"id", "firstName", "surname", "password", "roles"}).
		AddRow("1", "fname", "surname", "password", "{admin,user}")

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, firstName, surname, password, roles FROM user WHERE email = ?;`)).
		WithArgs(
			"test@example.com",
		).WillReturnRows(userMockRows)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, firstName, surname, password, roles FROM user WHERE email = ?;`)).
		WithArgs(
			"test2@example.com",
		).WillReturnError(sql.ErrNoRows)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := tc.ur.GetByEmail(tc.email)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\n want: %s\nhave: %s\n", tc.expectedErr, err.Error())
			}

			if !reflect.DeepEqual(u, domain.User{}) {
				iBytes, err := json.Marshal(u)
				if err != nil {
					t.Error(err)
				}

				if tc.expectedResult != string(iBytes) {
					t.Errorf("unexpected result \n want %s\nhave: %s\n", tc.expectedResult, string(iBytes))
				}
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type testcase struct {
		name           string
		email          string
		password       string
		id             uint64
		ur             ports.UserRepository
		expectedErr    error
		expectedResult string
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockRepo := NewUserRepository(db)

	testcases := []testcase{
		{
			name:        "get user success",
			email:       "test@example.com",
			password:    "password",
			id:          1,
			ur:          mockRepo,
			expectedErr: nil,
		},
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE user SET password TO ? WHERE email = ? AND id = ?;`)).
		WithArgs(
			"test@example.com",
			"password",
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ur.ChangePassword(tc.id, tc.email, tc.password)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\n want: %s\nhave: %s\n", tc.expectedErr, err.Error())
			}
		})
	}
}
