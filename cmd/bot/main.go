package main

import (
	"fmt"

	"github.com/CrabRus/LiveStats/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Бот %s готов к подключению на канал %s\n", cfg.Bot.BotName, cfg.Bot.Channel)
}
