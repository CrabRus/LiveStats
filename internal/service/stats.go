package service

import (
	"context"
	"log"

	"github.com/CrabRus/LiveStats/internal/domain"
)

type StatsService struct {
	s       domain.StreamsRepository
	w       domain.WordRepository
	u       domain.UserRepository
	channel string
}

func NewStatsService(s domain.StreamsRepository, w domain.WordRepository, u domain.UserRepository, channel string) *StatsService {
	return &StatsService{
		s:       s,
		w:       w,
		u:       u,
		channel: channel,
	}
}

func (s *StatsService) ProcessPeriodStats(ctx context.Context, stats *domain.PeriodStats) error {
	streamID, err := s.s.GetActiveStreamID(ctx, s.channel)
	if err != nil {
		log.Printf("Стрим для канала %s сейчас не активен в БД, пропускаем сбор за период %s",
			s.channel, stats.StartedAt.Format("15:04"))
		return nil
	}
	stats.StreamID = streamID
	if err := s.w.SaveStats(ctx, stats); err != nil {
		return err
	}

	log.Printf("Успешно сохранена статистика слов для стрима %s за период %s",
		streamID, stats.StartedAt.Format("15:04"))
	return nil
}
