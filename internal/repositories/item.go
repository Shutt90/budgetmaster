package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type itemRepository struct {
	*sql.DB
	clock ports.Clock
}

var (
	ErrNotFound      = errors.New("not found")
	ErrUnprocessable = errors.New("unprocessable entity")
)

func NewItemRepository(db *sql.DB, clock ports.Clock) *itemRepository {
	if db == nil || clock == nil {
		log.Error(fmt.Printf("db:%#v\nclock:%#v\n", db, clock))
		panic("unable to start new item repository")
	}
	return &itemRepository{db, clock}
}

func (db *itemRepository) CreateItemTable() error {
	db.DB.Begin()

	queryBytes, err := os.ReadFile("internal/migrations/item_table_schema.sql")
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = db.Exec(string(queryBytes))
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (db *itemRepository) Create(i domain.Item) error {
	query := "INSERT INTO item (name, description, location, cost, isRecurring) VALUES (?, ?, ?, ?, ?);"

	_, err := db.Exec(query, i.Name, i.Description, i.Location, i.Cost, i.IsRecurring)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (db *itemRepository) Get(id uint64) (domain.Item, error) {
	db.Begin()
	i := domain.Item{}

	row := db.QueryRow("SELECT * FROM item WHERE id = ?;", id)
	if row.Err() == sql.ErrNoRows {
		log.Error(row.Err())
		return domain.Item{}, ErrNotFound
	}

	if err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Location,
		&i.Cost,
		&i.IsRecurring,
		&i.RemovedOccuringAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	); err != nil {
		return domain.Item{}, err
	}

	return i, nil
}

func (db *itemRepository) GetMonthlyItems(month string, year int) ([]domain.Item, error) {
	items := []domain.Item{}

	rows, err := db.Query("SELECT * FROM item WHERE MONTH(createdAt) = ? AND YEAR(createdAt) = ?;", month, year)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error(err)
			return nil, ErrNotFound
		}
		log.Error(err)
		return []domain.Item{}, ErrUnprocessable
	}

	for rows.Next() {
		i := domain.Item{}
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Location,
			&i.Cost,
			&i.IsRecurring,
			&i.RemovedOccuringAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return []domain.Item{}, err
		}

		items = append(items, i)
	}

	return items, nil
}

func (db *itemRepository) SwitchRecurringPayments(id uint64, isRecurring bool) error {
	db.Begin()

	row := db.QueryRow("SELECT isRecurring FROM item WHERE id = ?;", id)

	i := domain.Item{}
	if err := row.Scan(&i.IsRecurring); err != nil {
		log.Error(err)
		return ErrNotFound
	}

	if i.IsRecurring == isRecurring {
		log.Error(ErrUnprocessable.Error())
		return ErrUnprocessable
	}

	if i.IsRecurring == true {
		_, err := db.Exec("UPDATE item SET isRecurring = ?, removedRecurringAt = ? WHERE id = ?;", isRecurring, db.clock.Now(), id)
		if err != nil {
			return err
		}

		return nil
	}

	_, err := db.Exec("UPDATE item SET isRecurring = ? WHERE id = ?;", isRecurring, id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
