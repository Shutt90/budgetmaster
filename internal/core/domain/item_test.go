package domain

import (
	"testing"
)

func TestNewItem(t *testing.T) {
	t.Run("checks item created", func(t *testing.T) {
		i := NewItem(
			"name",
			"description",
			"location",
			"month",
			2024,
			100,
			true,
		)

		if i == nil {
			t.Error("new item not created")
		}

		if i.Name != "name" {
			t.Error("incorrect item name")
		}

		if i.Description != "description" {
			t.Error("incorrect item description")
		}

		if i.Location != "location" {
			t.Error("incorrect location")
		}

		if i.Month != "month" {
			t.Error("incorrect month")
		}

		if i.Year != 2024 {
			t.Error("incorrect year")
		}

		if i.Cost != 100 {
			t.Error("incorrect item cost")
		}

		if i.IsRecurring != true {
			t.Error("incorrect recurring")
		}

	})
}
