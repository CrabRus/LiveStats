package service

import (
	"context"
	"log"
	"time"
)

type PeriodStats struct {
	StreamID  string
	StartedAt time.Time
	Words     map[string]int
}

// Описываем, какие методы нам нужны от репозитория
type WordRepository interface {
	GetActiveStreamID(ctx context.Context, channel string) (string, error)
	SaveStats(ctx context.Context, stats *PeriodStats) error // обновили тип
}

type StatsService struct {
	repo    WordRepository
	channel string // имя канала берем из конфига, чтобы знать, для кого искать стрим
}

func NewStatsService(repo WordRepository, channel string) *StatsService {
	return &StatsService{
		repo:    repo,
		channel: channel,
	}
}

// ProcessPeriodStats обрабатывает 5-минутку: находит стрим и отправляет в репозиторий
func (s *StatsService) ProcessPeriodStats(ctx context.Context, stats *PeriodStats) error {
	// 1. Пытаемся получить ID активного стрима
	streamID, err := s.repo.GetActiveStreamID(ctx, s.channel)
	if err != nil {
		// Если стрим не найден в БД (стример оффлайн), мы просто пропускаем сохранение,
		// чтобы не копить в базе "мусорные" слова, сказанные при выключенном стриме.
		log.Printf("Стрим для канала %s сейчас не активен в БД, пропускаем сбор за период %s",
			s.channel, stats.StartedAt.Format("15:04"))
		return nil
	}

	// 2. Устанавливаем ID стрима, который мы только что нашли
	stats.StreamID = streamID

	// 3. Отправляем в базу на сохранение (UPSERT)
	if err := s.repo.SaveStats(ctx, stats); err != nil {
		return err
	}

	log.Printf("Успешно сохранена статистика слов для стрима %s за период %s",
		streamID, stats.StartedAt.Format("15:04"))
	return nil
}
