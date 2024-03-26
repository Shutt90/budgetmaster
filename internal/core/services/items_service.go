package services

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

func NewItem(name, description string, cost uint64, expiry time.Time) *Item {
	return &Item{
		name:        name,
		description: description,
		cost:        cost,
		expiry:      expiry,
	}
}

func (srv *Item) Show() []Item {

}
