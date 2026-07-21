package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server    ServerConfig
	Bot       BotConfig
	Database  DatabaseConfig
	TickerMin int
}

type ServerConfig struct {
	Host string
	Port string
	Env  string
}

type BotConfig struct {
	Token   string
	BotName string
	Channel string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Bot: BotConfig{
			Token:   getEnv("TWITCH_TOKEN", "---"),
			BotName: getEnv("TWITCH_BOT_NAME", "---"),
			Channel: getEnv("TWITCH_CHANNEL", "---"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "livestats_user"),
			Password: getEnv("DB_PASSWORD", "12345678"),
			DBName:   getEnv("DB_NAME", "livestats_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		TickerMin: getEnvAsInt("TICKER_MIN", 1),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		// Если в .env передали не число, фолбэчимся на дефолт
		return defaultValue
	}

	return value
}
