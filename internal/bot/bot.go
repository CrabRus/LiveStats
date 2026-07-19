package bot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/CrabRus/LiveStats/internal/config"
	"github.com/CrabRus/LiveStats/internal/service"
	"github.com/gempir/go-twitch-irc/v4"
)

type Bot struct {
	cfg          *config.Config
	client       *twitch.Client
	statsService *service.StatsService // Используем указатель на сервис напрямую
	mu           sync.Mutex
	stats        *service.PeriodStats // Используем структуру из пакета service
}

func New(cfg *config.Config, svc *service.StatsService) *Bot {
	client := twitch.NewClient(cfg.Bot.BotName, cfg.Bot.Token)

	b := &Bot{
		cfg:          cfg,
		client:       client,
		statsService: svc,
		stats: &service.PeriodStats{ // Используем структуру из пакета service
			StreamID:  "",
			StartedAt: time.Now(),
			Words:     make(map[string]int),
		},
	}

	b.setupHandlers()

	return b
}

func (b *Bot) setupHandlers() {
	b.client.OnConnect(b._OnConnect)
	b.client.OnPrivateMessage(b._OnPrivateMessage)
}

func (b *Bot) _OnConnect() {
	fmt.Printf("Bot %s successfully connected to %s!\n", b.cfg.Bot.BotName, b.cfg.Bot.Channel)
}

func (b *Bot) _OnPrivateMessage(message twitch.PrivateMessage) {
	words := CleanAndSplit(message.Message)
	if len(words) == 0 {
		return
	}

	b.mu.Lock()
	for _, word := range words {
		b.stats.Words[word]++
	}
	b.mu.Unlock()
}

func (b *Bot) Start() error {
	b.client.Join(b.cfg.Bot.Channel)
	go b.startTicker()
	return b.client.Connect()
}

func (b *Bot) startTicker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()

		currentStats := b.stats

		b.stats = &service.PeriodStats{
			StreamID:  "",
			StartedAt: time.Now(),
			Words:     make(map[string]int),
		}
		b.mu.Unlock()
		b.sendToService(currentStats)
	}
}

func (b *Bot) sendToService(stats *service.PeriodStats) {
	// Создаем контекст с таймаутом на отправку в сервис/БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Передаем собранные данные в слой бизнес-логики
	if err := b.statsService.ProcessPeriodStats(ctx, stats); err != nil {
		log.Printf("Ошибка при обработке статистики сервисом: %v", err)
	}
}
