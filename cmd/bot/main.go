package main

import (
	"log"

	"github.com/CrabRus/LiveStats/internal/bot"
	"github.com/CrabRus/LiveStats/internal/config"
	"github.com/CrabRus/LiveStats/internal/db"
	"github.com/CrabRus/LiveStats/internal/repository/postgres"
	"github.com/CrabRus/LiveStats/internal/service"
)

func main() {
	cfg := config.Load()
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	streamsRepo := postgres.NewStreamsRepository(database)
	wordRepo := postgres.NewWordRepository(database)
	userRepo := postgres.NewUserRepository(database)

	statsService := service.NewStatsService(streamsRepo, wordRepo, userRepo, cfg.Bot.Channel)

	twitchBot := bot.New(cfg, statsService)

	log.Printf("Starting bot for channel %s...", cfg.Bot.Channel)
	if err := twitchBot.Start(); err != nil {
		log.Fatalf("Bot runtime error: %v", err)
	}
}
