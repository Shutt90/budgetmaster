package domain

import (
	"database/sql"
)

type Item struct {
	ID                uint64        `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Location          string        `json:"location"`
	Cost              uint64        `json:"cost"`
	Month             string        `json:"month"`
	Year              uint16        `json:"year"`
	IsRecurring       bool          `json:"isRecurring"`
	RemovedOccuringAt *sql.NullTime `json:"removedOccuringAt,omitempty"`
	CreatedAt         *sql.NullTime `json:"createdAt"`
	UpdatedAt         *sql.NullTime `json:"updatedAt,omitempty"`
}

func NewItem(name, desc, loc, month string, year uint16, cost uint64, isRecurring bool) *Item {
	return &Item{
		Name:        name,
		Description: desc,
		Location:    loc,
		Cost:        cost,
		Month:       month,
		Year:        year,
		IsRecurring: isRecurring,
	}
}
