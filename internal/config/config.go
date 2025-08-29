package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	DBDsn     string
	JWTSecret string
	Env       string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		host := get("DB_HOST", "localhost")
		port := get("DB_PORT", "5432")
		user := get("DB_USER", "postgres")
		pass := get("DB_PASSWORD", "123")
		name := get("DB_NAME", "logistic_db")
		dsn = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + name + "?sslmode=disable"
	}

	return &Config{
		AppPort:   get("APP_PORT", "8081"),
		DBDsn:     dsn,
		JWTSecret: get("JWT_SECRET", "eeee1234qq11"),
		Env:       get("ENV", "dev"),
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
