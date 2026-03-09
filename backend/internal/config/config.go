package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	JWTSecret  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBName:     getEnv("DB_NAME", "open_talk.db"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		JWTSecret:  loadJWTSecret(),
	}
}

func (c *Config) IsSQLite() bool {
	return c.DBHost == ""
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadJWTSecret() string {
	if value, ok := os.LookupEnv("JWT_SECRET"); ok && value != "" {
		return value
	}
	// Keep fallback for development convenience, but warn loudly.
	log.Println("WARNING: JWT_SECRET is not set. Using insecure default. Set JWT_SECRET in production.")
	return "default-secret-change-me"
}
