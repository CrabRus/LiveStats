package postgres

import (
	"context"
	"time"

	"github.com/CrabRus/LiveStats/internal/db"
	"github.com/CrabRus/LiveStats/internal/domain"
)

type WordRepository struct {
	db *db.DB
}

func NewWordRepository(database *db.DB) domain.WordRepository {
	return &WordRepository{db: database}
}

func (w *WordRepository) SaveStats(ctx context.Context, stats *domain.PeriodStats) error {
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

	_, err := w.db.NamedExecContext(ctx, query, rows)
	return err
}
