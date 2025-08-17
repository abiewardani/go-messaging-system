package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver
)

type PostgresDB struct {
    *sql.DB
}

func NewPostgresDB(dataSourceName string) (*PostgresDB, error) {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    log.Println("Connected to PostgreSQL database")
    return &PostgresDB{db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
    return p.DB.Close()
}