package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
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

func (db *database) Get(id uint64) (*domain.Item, error) {
	db.Begin()
	i := domain.Item{}

	row := db.QueryRow("SELECT * FROM items WHERE id = ?", id)
	if row.Err() == sql.ErrNoRows {
		return nil, fmt.Errorf("%s", ErrNotFound)
	}

	if err := row.Scan(&i); err != nil {
		return nil, fmt.Errorf(ErrUnprocessable)
	}

	return &i, nil
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
