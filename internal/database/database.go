package database

import (
	"avito-task/internal/config"
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectToDatabase() (*sql.DB, error) {
	conn := config.Get().DatabaseDSN
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
