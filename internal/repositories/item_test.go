package repositories

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
	"github.com/Shutt90/budgetmaster/internal/core/services"
)

type testcase struct {
	name        string
	item        func() *domain.Item
	ir          ports.ItemRepository
	expectedErr error
}

func TestCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()

	mockClock := services.NewMockClock()

	mockRepo := NewItemRepository(db, mockClock)

	testcases := []testcase{
		{
			name: "test new user created success",
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
