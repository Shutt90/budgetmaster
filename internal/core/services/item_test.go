package services

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
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

	mockService := NewItemService(repositories.NewItemRepository(db, NewMockClock()), mockClock)

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
