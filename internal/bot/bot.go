package bot

import (
	"fmt"
	"sync"

	"github.com/CrabRus/LiveStats/internal/config"
	"github.com/gempir/go-twitch-irc/v4"
)

var (
	WordCount = make(map[string]int)
)

type Bot struct {
	cfg       *config.Config
	client    *twitch.Client
	mu        sync.Mutex
	wordCount map[string]int
}

func New(cfg *config.Config) *Bot {
	client := twitch.NewClient(cfg.Bot.BotName, cfg.Bot.Token)

	b := &Bot{
		cfg:       cfg,
		client:    client,
		wordCount: make(map[string]int),
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
		b.wordCount[word]++
	}
	b.mu.Unlock()
}

func (b *Bot) Start() error {
	b.client.Join(b.cfg.Bot.Channel)
	return b.client.Connect()
}
