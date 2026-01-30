package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// connection pool (AMAN untuk pooler)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	log.Println("✅ Database connected successfully")
	return db, nil
}
