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
					2024,
					100,
					false,
				)
			},
			ir:          mockRepo,
			expectedErr: fmt.Errorf(""),
		},
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO item (name, description, location, cost, month, year, isMonthly) VALUES (?, ?, ?, ?, ?, ?);`)).
		WithArgs(
			"testName",
			"testDesc",
			"testLoc",
			100,
			"testMonth",
			2024,
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
			expectedResult: `{"ID":1,"Name":"testName","Description":"testDesc","Location":"testLoc","Cost":100,"Month":"April","Year":2024,"IsRecurring":true,"RemovedOccuringAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"CreatedAt":{"Time":"2024-03-27T16:26:00Z","Valid":true},"UpdatedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`,
		},
		{
			name:           "user not found",
			id:             2,
			ir:             mockRepo,
			expectedErr:    ErrNotFound,
			expectedResult: `{}`,
		},
	}

	itemMockRows := sqlmock.NewRows([]string{"id", "name", "description", "location", "cost", "month", "year", "isRecurring", "removedRecurringAt", "createdAt", "updatedAt"}).
		AddRow("1", "testName", "testDesc", "testLoc", 100, "April", 2024, "1", sql.NullTime{}, mockRepo.clock.Now(), sql.NullTime{})

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM item WHERE id = ?;`)).
		WithArgs(
			1,
		).WillReturnRows(itemMockRows)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM item WHERE id = ?;`)).
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

func TestGetMonthlyItems(t *testing.T) {
	type testcase struct {
		name           string
		month          string
		year           uint16
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
			name:           "get monthly for period of jan",
			year:           2024,
			month:          "January",
			ir:             mockRepo,
			expectedErr:    nil,
			expectedResult: `[{"ID":1,"Name":"testName","Description":"testDesc","Location":"testLoc","Cost":100,"Month":"January","Year":2024,"IsRecurring":true,"RemovedOccuringAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"CreatedAt":{"Time":"2024-03-27T16:26:00Z","Valid":true},"UpdatedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}},{"ID":2,"Name":"testName2","Description":"testDesc2","Location":"testLoc2","Cost":200,"Month":"January","Year":2024,"IsRecurring":false,"RemovedOccuringAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"CreatedAt":{"Time":"2024-03-27T16:26:00Z","Valid":true},"UpdatedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}]`,
		},
	}

	itemMockRows := sqlmock.NewRows([]string{"id", "name", "description", "location", "cost", "month", "year", "isRecurring", "removedRecurringAt", "createdAt", "updatedAt"}).
		AddRow(1, "testName", "testDesc", "testLoc", 100, "January", 2024, "1", sql.NullTime{}, mockRepo.clock.Now(), sql.NullTime{}).
		AddRow(2, "testName2", "testDesc2", "testLoc2", 200, "January", 2024, "0", sql.NullTime{}, mockRepo.clock.Now(), sql.NullTime{})

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM item WHERE month = ? AND year = ?;`)).
		WithArgs(
			"January",
			2024,
		).WillReturnRows(itemMockRows)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			i, err := tc.ir.GetMonthlyItems("January", 2024)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\n want: %s\nhave: %s\n", tc.expectedErr, err.Error())
			}

			iBytes, err := json.Marshal(i)
			if err != nil {
				t.Error(err)
			}

			if tc.expectedResult != string(iBytes) {
				t.Errorf("unexpected result \n want %s\nhave: %s\n", tc.expectedResult, string(iBytes))
			}
		})
	}
}

func TestSwitchRecurringPayments(t *testing.T) {
	type testcase struct {
		name        string
		id          uint64
		isRecurring bool
		ir          ports.ItemRepository
		expectedErr error
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mockClock := services.NewMockClock()
	mockRepo := NewItemRepository(db, mockClock)

	testcases := []testcase{
		{
			name:        "switch from recurring to not",
			id:          1,
			isRecurring: false,
			ir:          mockRepo,
			expectedErr: nil,
		},
		{
			name:        "switch from not recurring to recurring",
			id:          2,
			isRecurring: true,
			ir:          mockRepo,
			expectedErr: nil,
		},
	}

	itemMockRows1 := sqlmock.NewRows([]string{"isRecurring"}).
		AddRow(1)

	itemMockRows2 := sqlmock.NewRows([]string{"isRecurring"}).
		AddRow(0)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT isRecurring FROM item WHERE id = ?;`)).
		WithArgs(
			1,
		).WillReturnRows(itemMockRows1)

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE item SET isRecurring = ?, removedRecurringAt = ? WHERE id = ?;`)).
		WithArgs(
			false,
			mockClock.Now(),
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT isRecurring FROM item WHERE id = ?;`)).
		WithArgs(
			2,
		).WillReturnRows(itemMockRows2)

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE item SET isRecurring = ? WHERE id = ?;`)).
		WithArgs(
			true,
			2,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ir.SwitchRecurringPayments(tc.id, tc.isRecurring)
			if err != tc.expectedErr {
				t.Errorf("unexpected error\n want: %s\nhave: %s\n", tc.expectedErr, err.Error())
			}
		})
	}
}
