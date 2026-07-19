package postgres

import (
	"context"
	"time"

	"github.com/CrabRus/LiveStats/internal/db"
	"github.com/CrabRus/LiveStats/internal/service"
)

type WordRepository struct {
	db *db.DB
}

func NewWordRepository(database *db.DB) *WordRepository {
	return &WordRepository{db: database}
}

func (r *WordRepository) GetActiveStreamID(ctx context.Context, channel string) (string, error) {
	var streamID string
	query := `SELECT id FROM streams WHERE channel = $1 AND is_active = true LIMIT 1`

	err := r.db.GetContext(ctx, &streamID, query, channel)
	if err != nil {
		return "", err
	}

	return streamID, nil
}

func (r *WordRepository) SaveStats(ctx context.Context, stats *service.PeriodStats) error {
	if len(stats.Words) == 0 {
		return nil
	}

	query := `
		INSERT INTO word_stats (stream_id, time_frame, word, count)
		VALUES (:stream_id, :time_frame, :word, :count)
		ON CONFLICT (stream_id, time_frame, word) 
		DO UPDATE SET count = word_stats.count + EXCLUDED.count;
	`
	type wordRow struct {
		StreamID  string    `db:"stream_id"`
		TimeFrame time.Time `db:"time_frame"`
		Word      string    `db:"word"`
		Count     int       `db:"count"`
	}

	rows := make([]wordRow, 0, len(stats.Words))
	for word, count := range stats.Words {
		rows = append(rows, wordRow{
			StreamID:  stats.StreamID,
			TimeFrame: stats.StartedAt,
			Word:      word,
			Count:     count,
		})
	}

	_, err := r.db.NamedExecContext(ctx, query, rows)
	return err
}
