package repositories

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type database struct {
	db *sql.DB
}

func New(dsn string, token string) *database {
	url := fmt.Sprintf("%s?authToken=%s", dsn, token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()

	return &database{
		db: db,
	}
}
