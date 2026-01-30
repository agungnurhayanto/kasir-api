package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {

	dsn := connectionString + " prefer_simple_protocol=true"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("✅ Database connected successfully")
	return db, nil
}
