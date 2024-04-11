package domain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewItem(t *testing.T) {
	t.Run("checks item created", func(t *testing.T) {
		expected := Item{
			Name:        "name",
			Description: "description",
			Location:    "location",
			Month:       "month",
			Year:        2024,
			Cost:        100,
			IsRecurring: true,
		}

		i := NewItem(
			"name",
			"description",
			"location",
			"month",
			2024,
			100,
			true,
		)

		diff := cmp.Diff(i, expected)
		if diff != "" {
			t.Errorf("expected no diff but got %s", diff)
		}
	})
}
