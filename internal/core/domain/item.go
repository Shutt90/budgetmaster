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

type Clock struct{}

func (c Clock) Now() time.Time {
	return time.Now()
}

type ItemService struct {
	item  *Item
	clock clockIface
}

type clockIface interface {
	Now() time.Time
}

func NewItemService(i *Item, c clockIface) *ItemService {
	return &ItemService{
		item:  i,
		clock: c,
	}
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

func (i *ItemService) RemoveOccuring() {
	i.item.IsRecurring = false
	i.item.RemovedOccuringAt = i.clock.Now()
}
