package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
)

func NewDB(dbName string, auth string) *sql.DB {
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, auth)

	db, err := sql.Open("libsql", url)
	if err != nil {
		log.Error(err)
		os.Exit(1)

		return &sql.DB{}
	}

	_, err = db.Conn(context.Background())
	if err != nil {
		log.Error(err)
		return &sql.DB{}
	}

	log.Info("db connected")

	return db
}
