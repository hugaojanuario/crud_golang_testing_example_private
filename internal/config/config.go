package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost     string
	DBPort     string
	DbUser     string
	DbPassword string
	SslMode    string
}

func NewConfig() *Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Erro ao ler as variaveis de ambiente em '.env':", err)
	}

	return &Config{
		DbHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		SslMode:    os.Getenv("DB_SSLMODE"),
	}
}
