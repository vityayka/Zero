package models

import (
	"database/sql"
	"fmt"
)

func Open(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, err
	}
	return db, err
}

func DefaultPgConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "vityayka",
		Password: "azazaz1488",
		Database: "zero",
		Sslmode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Sslmode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User,
		cfg.Password, cfg.Database, cfg.Sslmode)
}
