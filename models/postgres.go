package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func Open(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, err
	}
	return db, err
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

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, fs fs.FS, dir string) error {
	goose.SetBaseFS(fs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
