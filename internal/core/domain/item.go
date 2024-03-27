package domain

import (
	"database/sql"
)

type Item struct {
	ID                uint64
	Name              string
	Description       string
	Location          string
	Cost              uint64
	Month             string
	IsRecurring       bool
	RemovedOccuringAt sql.NullTime
	CreatedAt         sql.NullTime
	UpdatedAt         sql.NullTime
}

func NewItem(name, desc, loc, month string, cost uint64, isRecurring bool) *Item {
	return &Item{
		Name:        name,
		Description: desc,
		Location:    loc,
		Cost:        cost,
		Month:       month,
		IsRecurring: isRecurring,
	}
}
