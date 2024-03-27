package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
	"github.com/Shutt90/budgetmaster/internal/core/services"
)

func TestCreate(t *testing.T) {
	type testcase struct {
		name        string
		item        func() *domain.Item
		ir          ports.ItemRepository
		expectedErr error
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockClock := services.NewMockClock()

	mockRepo := NewItemRepository(db, mockClock)

	testcases := []testcase{
		{
			name: "create new user success",
			item: func() *domain.Item {
				return domain.NewItem(
					"testName",
					"testDesc",
					"testLoc",
					"testMonth",
					100,
					false,
				)
			},
			ir:          mockRepo,
			expectedErr: fmt.Errorf(""),
		},
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO item (name, description, location, cost, month, isMonthly) VALUES (?, ?, ?, ?, ?, ?);`)).
		WithArgs(
			"testName",
			"testDesc",
			"testLoc",
			100,
			"testMonth",
			false,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.ir.Create(*tc.item()); err != nil {
				t.Errorf("err when adding new item, %s", err.Error())
			}
		})
	}
}

func TestGet(t *testing.T) {
	type testcase struct {
		name           string
		id             uint64
		ir             ports.ItemRepository
		expectedErr    error
		expectedResult string
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockClock := services.NewMockClock()
	mockRepo := NewItemRepository(db, mockClock)

	testcases := []testcase{
		{
			name:           "get user success",
			id:             1,
			ir:             mockRepo,
			expectedErr:    nil,
			expectedResult: `{"ID":1,"Name":"testName","Description":"testDesc","Location":"testLoc","Cost":100,"Month":"April","IsRecurring":true,"RemovedOccuringAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"CreatedAt":{"Time":"2024-03-27T16:26:00Z","Valid":true},"UpdatedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`,
		},
		{
			name:           "user not found",
			id:             2,
			ir:             mockRepo,
			expectedErr:    ErrNotFound,
			expectedResult: `{}`,
		},
	}

	itemMockRows := sqlmock.NewRows([]string{"id", "name", "description", "location", "cost", "month", "isRecurring", "removedRecurringAt", "createdAt", "updatedAt"}).
		AddRow("1", "testName", "testDesc", "testLoc", 100, "April", "1", sql.NullTime{}, mockRepo.clock.Now(), sql.NullTime{})

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM items WHERE id = ?;`)).
		WithArgs(
			1,
		).WillReturnRows(itemMockRows)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM items WHERE id = ?;`)).
		WithArgs(
			2,
		).WillReturnError(sql.ErrNoRows)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			i, err := tc.ir.Get(tc.id)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\n want: %s\nhave: %s\n", tc.expectedErr, err.Error())
			}

			if !reflect.DeepEqual(i, domain.Item{}) {
				iBytes, err := json.Marshal(i)
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
