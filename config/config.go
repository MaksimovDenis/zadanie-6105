package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress    string `default:"0.0.0.0:8080"`
	PostgresConn     string `default:"postgres://admin:admin@0.0.0.0:5432/avito?sslmode=disable"`
	PostgresJdbcURL  string `default:"jdbc:postgresql://0.0.0.0:5432/avito"`
	PostgresUsername string `default:"admin"`
	PostgresPassword string `default:"admin"`
	PostgresHost     string `default:"0.0.0.0"`
	PostgresPort     string `default:"5432"`
	PostgresDatabase string `default:"avito"`
}

func InitConfig() (*Config, error) {
	var cnf Config

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cnf = Config{
		ServerAddress:    os.Getenv("SERVER_ADDRESS"),
		PostgresConn:     os.Getenv("POSTGRES_CONN"),
		PostgresJdbcURL:  os.Getenv("POSTGRES_JDBC_URL"),
		PostgresUsername: os.Getenv("POSTGRES_USERNAME"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDatabase: os.Getenv("POSTGRES_DATABASE"),
	}

	return &cnf, nil
}
