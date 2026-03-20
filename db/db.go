package db

import (
	"fmt"
	"os"

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

	// Read SQL schema from file
	schemaBytes, err := os.ReadFile("db/db_sql.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	schema := string(schemaBytes)
	_, err = DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
