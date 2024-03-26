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

func New(dsn string, token string) *database {
	url := fmt.Sprintf("%s?authToken=%s", dsn, token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()

	return &database{db}
}

func (db *database) Get(id string) (*domain.Item, error) {
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
