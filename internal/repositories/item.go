package repositories

import (
	"database/sql"
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

const ErrNotFound = "not found"
const ErrUnprocessable = "unprocessable entity"

func NewDB(dsn string) *sql.DB {
	url := fmt.Sprintf("%s", dsn)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	return db
}

func NewItemRepository(db *sql.DB, clock ports.Clock) *itemRepository {
	return &itemRepository{db, clock}
}

func (db *itemRepository) CreateItemTable() error {
	db.DB.Begin()

	queryBytes, err := os.ReadFile("internal/migrations/item_table_schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(queryBytes))
	if err != nil {
		return err
	}

	return nil
}

func (db *itemRepository) Create(i domain.Item) error {
	query := "INSERT INTO item (name, description, location, cost, month, isMonthly) VALUES (?, ?, ?, ?, ?, ?);"

	_, err := db.Exec(query, i.Name, i.Description, i.Location, i.Cost, i.Month, i.IsRecurring)
	if err != nil {
		return err
	}

	return nil
}

func (db *itemRepository) Get(id uint64) (domain.Item, error) {
	db.Begin()
	i := domain.Item{}

	row := db.QueryRow("SELECT * FROM items WHERE id = ?;", id)
	if row.Err() == sql.ErrNoRows {
		return domain.Item{}, fmt.Errorf("%s", ErrNotFound)
	}

	if err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Location,
		&i.Cost,
		&i.Month,
		&i.IsRecurring,
		&i.RemovedOccuringAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	); err != nil {
		return domain.Item{}, fmt.Errorf(err.Error())
	}

	return i, nil
}

func (db *itemRepository) GetMonthlyItems(month string, year int) ([]domain.Item, error) {
	db.Begin()
	items := []domain.Item{}

	rows, err := db.Query("SELECT * FROM items WHERE month = ? AND year = ?;", month, year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s", ErrNotFound)
		}
		log.Error(err)
	}

	for rows.Next() {
		i := domain.Item{}
		if err := rows.Scan(i); err != nil {
			return []domain.Item{}, err
		}

		items = append(items, i)
	}

	return items, nil
}

func (db *itemRepository) SwitchRecurringPayments(id uint64, isRecurring bool) error {
	db.Begin()

	row := db.QueryRow("SELECT isRecurring FROM item WHERE id = ?")

	i := domain.Item{}
	if err := row.Scan(&i.IsRecurring); err != nil {
		return fmt.Errorf(ErrNotFound)
	}

	if i.IsRecurring == isRecurring {
		return fmt.Errorf(ErrUnprocessable)
	}

	if i.IsRecurring == true {
		_, err := db.Exec("UPDATE item SET isRecurring = ?, removedRecurringAt = ? WHERE id = ?;", isRecurring, db.clock.Now(), id)
		if err != nil {
			return err
		}
	}

	_, err := db.Exec("UPDATE item SET isRecurring = ? WHERE id = ?;", isRecurring, id)
	if err != nil {
		return err
	}

	return nil
}
