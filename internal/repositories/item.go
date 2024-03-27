package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/labstack/gommon/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type database struct {
	*sql.DB
}

const ErrNotFound = "not found"
const ErrUnprocessable = "unprocessable entity"

func New(dsn string) *database {
	url := fmt.Sprintf("%s", dsn)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	return &database{db}
}

func (db *database) CreateItemTable() error {
	db.Begin()

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

func (db *database) Create(i domain.Item) error {
	query := "INSERT INTO item (name, description, location, cost, month, isMonthly) VALUES (?, ?, ?, ?, ?, ?);"

	_, err := db.Exec(query, i.Name, i.Description, i.Location, i.Cost, i.Month, i.IsMonthly)
	if err != nil {
		return err
	}

	return nil
}

func (db *database) Get(id uint64) (domain.Item, error) {
	db.Begin()
	i := domain.Item{}

	row := db.QueryRow("SELECT * FROM items WHERE id = ?;", id)
	if row.Err() == sql.ErrNoRows {
		return domain.Item{}, fmt.Errorf("%s", ErrNotFound)
	}

	if err := row.Scan(&i); err != nil {
		return domain.Item{}, fmt.Errorf(ErrUnprocessable)
	}

	return i, nil
}

func (db *database) GetMonthlyItems(month string, year int) ([]domain.Item, error) {
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
