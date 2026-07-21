package domain

import (
	"context"
)

type WordRepository interface {
	SaveStats(ctx context.Context, stats *PeriodStats) error
}

type UserRepository interface {
	SaveStats(ctx context.Context, stats *PeriodStats) error
}

type StreamsRepository interface {
	GetActiveStreamID(ctx context.Context, channel string) (string, error)
}
