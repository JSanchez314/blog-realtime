package db

import (
	"github.com/jmoiron/sqlx"
)

func NewPostgres(pgURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", pgURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
