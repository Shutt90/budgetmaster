package domain

import "time"

type Item struct {
	ID                uint64
	Name              string
	Description       string
	Location          string
	Cost              uint64
	Month             string
	IsRecurring       bool
	RemovedOccuringAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
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
