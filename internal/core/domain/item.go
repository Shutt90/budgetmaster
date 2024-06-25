package domain

import (
	"database/sql"
	"time"
)

type Item struct {
	ID                uint64        `json:"id,omitempty"`
	Name              string        `json:"name"`
	Description       string        `json:"description,omitempty"`
	Location          string        `json:"location"`
	Cost              uint64        `json:"cost"`
	IsRecurring       bool          `json:"isRecurring"`
	RemovedOccuringAt *sql.NullTime `json:"removedOccuringAt,omitempty"`
	CreatedAt         *sql.NullTime `json:"createdAt,omitempty"`
	UpdatedAt         *sql.NullTime `json:"updatedAt,omitempty"`
}

func NewItem(name string, desc string, loc string, cost uint64, isRecurring bool) Item {
	return Item{
		Name:        name,
		Description: desc,
		Location:    loc,
		Cost:        cost,
		IsRecurring: isRecurring,
	}
}

func ParseCreatedAtString(datetime time.Time) string {
	return datetime.Format(time.DateTime)
}

func CostToFloat(cost int64) float64 {
	return float64(cost) / 100
}
