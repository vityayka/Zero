package main

import (
	"encoding/base64"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

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

func main() {
	bytes := make([]byte, 32)
	fmt.Printf("original bytes: %v\n", bytes)
	// rand.Read(bytes)
	// fmt.Printf("rand bytes: %v\n", bytes)
	str := base64.URLEncoding.EncodeToString(bytes)
	fmt.Printf("base64: %v", str)
}
