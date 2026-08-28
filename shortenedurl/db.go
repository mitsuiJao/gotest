package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB() *sql.DB {
	dsn := "postgres://urlshortener:password@localhost:5432/urlshortener?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err = migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	fmt.Println("db connected")
	return db
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			code TEXT PRIMARY KEY,
			long_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		)
	`) 
	return err
}
