package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBURL      string
	RabbitURL  string
	JWTSecret  string
	Port       string
}

func Load() *Config {
	// Load .env from current dir, fall back to ../.env (root) for local dev
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load("../.env")
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "forte"),
		DBPassword: getEnv("DB_PASSWORD", "forte123"),
		DBName:     getEnv("DB_NAME", "forte_commerce"),
		RabbitURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JWTSecret:  getEnv("JWT_SECRET", "forte-secret-key"),
		Port:       getEnv("PORT", "8080"),
	}

	// Construct DBURL
	cfg.DBURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
