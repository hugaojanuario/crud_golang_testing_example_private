package postgre

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	DbHost     string
	DBPort     string
	DbUser     string
	DbPassword string
	DbName     string
	SslMode    string
}

func NewConn(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DbHost, cfg.DBPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.SslMode)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("Error ao conectar com o banco de dados: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Erro ao estabelecer conexão com o banco de dados: %w", err)
	}

	return db, nil
}
