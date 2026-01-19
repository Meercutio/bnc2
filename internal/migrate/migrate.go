package migrate

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Up applies all pending migrations.
//
// It returns an error (no log.Fatal) so the caller can decide how to handle it.
func Up(dbURL, dir string, log *slog.Logger) error {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("migrations: open db: %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Error("database close error")
		}
	}(db)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrations: set dialect: %w", err)
	}

	if log != nil {
		log.Info("running database migrations", "dir", dir)
	}
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("migrations: goose up: %w", err)
	}
	if log != nil {
		log.Info("database migrations applied")
	}
	return nil
}
