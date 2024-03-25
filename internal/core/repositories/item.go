package repositories

import "github.com/Shutt90/budgetmaster/repositories"

type itemRepository struct {
	db repositories.Database
}

func NewItemRepository(db repositories.Database) *itemRepository {
	return &itemRepository{
		db: db,
	}
}
