package main

import (
	"log"

	"github.com/CrabRus/LiveStats/internal/bot"
	"github.com/CrabRus/LiveStats/internal/config"
)

func main() {
	cfg := config.Load()
	twitchBot := bot.New(cfg)

	log.Println("Запуск Twitch бота...")
	if err := twitchBot.Start(); err != nil {
		log.Fatalf("Ошибка при работе бота: %v", err)
	}
}
