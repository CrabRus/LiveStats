package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/CrabRus/LiveStats/internal/config"
	"github.com/gempir/go-twitch-irc/v4"
)

type PeriodStats struct {
	StreamID  string
	StartedAt time.Time
	Words     map[string]int
}

type Bot struct {
	cfg    *config.Config
	client *twitch.Client
	mu     sync.Mutex
	stats  *PeriodStats
}

func New(cfg *config.Config) *Bot {
	client := twitch.NewClient(cfg.Bot.BotName, cfg.Bot.Token)

	b := &Bot{
		cfg:    cfg,
		client: client,
		stats: &PeriodStats{
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

		b.stats = &PeriodStats{
			StreamID:  "",
			StartedAt: time.Now(),
			Words:     make(map[string]int),
		}
		b.mu.Unlock()
		b.sendToService(currentStats)
	}
}

func (b *Bot) sendToService(stats *PeriodStats) {

}
