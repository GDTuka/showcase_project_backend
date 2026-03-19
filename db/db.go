package db

import (
	"fmt"

	"showcase_project/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB(cfg *config.Config) error {
	var err error
	DB, err = sqlx.Connect("sqlite3", cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Create tables if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE NOT NULL,
		phone TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
