package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
	"github.com/Shutt90/budgetmaster/internal/core/services"
	"github.com/Shutt90/budgetmaster/internal/repositories"
)

func TestCreate(t *testing.T) {
	type testcase struct {
		name        string
		item        func() domain.Item
		ir          ports.ItemRepository
		expectedErr error
	}

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockClock := NewMockClock()

	// remove mock clocks into single struct
	mockService := NewItemService(repositories.NewItemRepository(db, mockClock), mockClock)

	testcases := []testcase{
		{
			name: "create new user success",
			item: func() domain.Item {
				return domain.NewItem(
					"testName",
					"testDesc",
					"testLoc",
					100,
					false,
				)
			},
			ir:          mockService.itemRepository,
			expectedErr: fmt.Errorf(""),
		},
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO item (name, description, location, cost, isRecurring) VALUES (?, ?, ?, ?, ?);`)).
		WithArgs(
			"testName",
			"testDesc",
			"testLoc",
			100,
			false,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.ir.Create(tc.item()); err != nil {
				t.Errorf("err when adding new item, %s", err.Error())
			}
		})
	}
}

func TestGetDefaultMonthlyItems(t *testing.T) {
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

	// remove mock clocks into single struct
	mockService := NewItemService(repositories.NewItemRepository(db, mockClock), mockClock)

	testcases := []testcase{
		{
			name:           "get monthly for period of jan",
			ir:             mockService.itemRepository,
			expectedErr:    nil,
			expectedResult: `[{"id":1,"name":"testName","description":"testDesc","location":"testLoc","cost":100,"isRecurring":true,"createdAt":{"Time":"2024-01-27T16:26:00Z","Valid":true},"updatedAt":{"Time":"2024-03-27T16:26:00Z","Valid":true}},{"id":2,"name":"testName2","description":"testDesc2","location":"testLoc2","cost":200,"isRecurring":false,"createdAt":{"Time":"2024-01-27T16:26:00Z","Valid":true},"updatedAt":{"Time":"2024-03-27T16:26:00Z","Valid":true}}]`,
		},
	}

	itemMockRows := sqlmock.NewRows([]string{"id", "name", "description", "location", "cost", "isRecurring", "removedRecurringAt", "createdAt", "updatedAt"}).
		AddRow(1, "testName", "testDesc", "testLoc", 100, "1", sql.NullTime{}, mockService.clock.Jan(), mockService.clock.Now()).
		AddRow(2, "testName2", "testDesc2", "testLoc2", 200, "0", sql.NullTime{}, mockService.clock.Jan(), mockService.clock.Now())

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM item WHERE strftime('%m', createdAt) = ? AND strftime('%Y', createdAt) = ?;`)).
		WithArgs(
			1,
			2024,
		).WillReturnRows(itemMockRows)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			i, err := tc.ir.GetMonthlyItems(1, 2024)
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
