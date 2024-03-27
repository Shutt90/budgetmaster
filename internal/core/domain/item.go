package domain

import "time"

type Item struct {
	ID          uint64
	Name        string
	Description string
	Location    string
	Cost        uint64
	Month       string
	IsMonthly   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
