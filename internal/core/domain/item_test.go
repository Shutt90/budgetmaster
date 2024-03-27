package domain

import (
	"testing"
)

func TestChangedRecurring(t *testing.T) {
	t.Run("checks item created", func(t *testing.T) {
		i := NewItem(
			"testItem",
			"testDescriptions",
			"testLocation",
			"testMonth",
			100,
			true,
		)

		if i == nil {
			t.Error("new item not created")
		}
	})
}
