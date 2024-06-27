package services

import "time"

type Clock struct{}

func NewClock() *Clock {
	return &Clock{}
}

func (c *Clock) Now() time.Time {
	return time.Now()
}

func (c *Clock) Jan() time.Time {
	return time.Date(2024, 01, 27, 16, 26, 0, 0, time.UTC)
}
