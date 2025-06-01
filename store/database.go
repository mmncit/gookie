package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {

	dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = fmt.Sprintf(
            "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
            os.Getenv("DB_HOST"),
            os.Getenv("DB_USER"),
            os.Getenv("DB_PASSWORD"),
            os.Getenv("DB_NAME"),
            os.Getenv("DB_PORT"),
            os.Getenv("DB_SSLMODE"),
        )
    }
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("connected to database...")
	return db, nil
}