package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func runMigrations(dbURL string) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("migrations: open db: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("migrations: set dialect: %v", err)
	}

	log.Println("running database migrations...")

	if err := goose.Up(db, "app/migrations/"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	log.Println("database migrations applied successfully")
}
