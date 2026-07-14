package config

import "os"

type Config struct {
	Server   ServerConfig
	Bot      BotConfig
	Database DatabaseConfig
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
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
