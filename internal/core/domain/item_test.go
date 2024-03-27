package domain

import (
	"testing"
	"time"
)

type mockClock struct {
	time.Time
}

var expectedTime = time.Date(2024, 3, 27, 15, 06, 30, 0, time.UTC)

func (mc mockClock) Now() time.Time {
	return expectedTime
}

func TestChangedRecurring(t *testing.T) {
	mockClock := mockClock{}

	i := NewItem(
		"testItem",
		"testDescriptions",
		"testLocation",
		"testMonth",
		100,
		true,
	)
	is := NewItemService(i, mockClock)

	t.Run("checks recurring at has changed and time added", func(t *testing.T) {
		is.RemoveOccuring()

		if is.item.IsRecurring == true {
			t.Errorf("expected %t, got %t", true, is.item.IsRecurring)
		}

		if is.item.RemovedOccuringAt.String() != expectedTime.String() {
			t.Errorf("expected %s, got %s", expectedTime, is.item.RemovedOccuringAt)
		}
	})

}
