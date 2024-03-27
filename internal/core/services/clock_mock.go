package services

import "time"

type mockClock struct{}

func NewMockClock() mockClock {
	return mockClock{}
}

func (mc mockClock) Now() time.Time {
	return time.Date(2024, 03, 27, 16, 26, 0, 0, time.UTC)
}
