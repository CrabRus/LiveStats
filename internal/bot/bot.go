package bot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/CrabRus/LiveStats/internal/config"
	"github.com/CrabRus/LiveStats/internal/domain"
	"github.com/CrabRus/LiveStats/internal/service"
	"github.com/gempir/go-twitch-irc/v4"
)

type Bot struct {
	cfg          *config.Config
	client       *twitch.Client
	statsService *service.StatsService
	mu           sync.Mutex
	stats        *domain.PeriodStats
}

func New(cfg *config.Config, svc *service.StatsService) *Bot {
	client := twitch.NewClient(cfg.Bot.BotName, cfg.Bot.Token)

	b := &Bot{
		cfg:          cfg,
		client:       client,
		statsService: svc,
		stats: &domain.PeriodStats{
			StreamID:  "",
			StartedAt: time.Now(),
			Words:     make(map[string]int),
			UserStats: make(map[string]int),
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
	defer b.mu.Unlock()

	if b.stats.Words == nil {
		b.stats.Words = make(map[string]int)
	}
	if b.stats.UserStats == nil {
		b.stats.UserStats = make(map[string]int)
	}

	for _, word := range words {
		b.stats.Words[word]++
	}

	username := message.User.Name
	b.stats.UserStats[username]++
}

func (b *Bot) Start() error {
	b.client.Join(b.cfg.Bot.Channel)
	go b.startTicker()
	return b.client.Connect()
}

func (b *Bot) startTicker() {
	ticker := time.NewTicker(time.Duration(b.cfg.TickerMin) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()

		currentStats := b.stats

		b.stats = &domain.PeriodStats{
			StreamID:  "",
			StartedAt: time.Now(),
			Words:     make(map[string]int),
		}
		b.mu.Unlock()
		b.sendToService(currentStats)
	}
}

func (b *Bot) sendToService(stats *domain.PeriodStats) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := b.statsService.ProcessPeriodStats(ctx, stats); err != nil {
		log.Printf("Ошибка при обработке статистики сервисом: %v", err)
	}
}
