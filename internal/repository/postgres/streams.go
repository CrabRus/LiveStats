package postgres

import (
	"context"

	"github.com/CrabRus/LiveStats/internal/db"
	"github.com/CrabRus/LiveStats/internal/domain"
)

type StreamsRepository struct {
	db *db.DB
}

func NewStreamsRepository(database *db.DB) domain.StreamsRepository {
	return &StreamsRepository{db: database}
}

func (s *StreamsRepository) GetActiveStreamID(ctx context.Context, channel string) (string, error) {
	var streamID string
	query := `SELECT id FROM streams WHERE channel = $1 AND is_active = true LIMIT 1`

	err := s.db.GetContext(ctx, &streamID, query, channel)
	if err != nil {
		return "", err
	}

	return streamID, nil
}
