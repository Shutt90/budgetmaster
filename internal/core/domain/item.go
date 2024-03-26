package domain

import "time"

type Item struct {
	id          uint64
	name        string
	description string
	cost        uint64
	expiry      time.Time
	createdAt   time.Time
	updatedAt   time.Time
}
