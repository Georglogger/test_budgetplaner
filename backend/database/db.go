package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func Connect(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Öffnen der Datenbankverbindung: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("fehler beim Verbinden zur Datenbank: %w", err)
	}

	log.Println("Datenbankverbindung erfolgreich hergestellt")
	return db, nil
}
