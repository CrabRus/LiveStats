package postgres

import (
	"context"
	"time"

	"github.com/CrabRus/LiveStats/internal/db"
	"github.com/CrabRus/LiveStats/internal/domain"
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(database *db.DB) domain.UserRepository {
	return &UserRepository{db: database}
}

func (u *UserRepository) SaveStats(ctx context.Context, stats *domain.PeriodStats) error {
	if len(stats.UserStats) == 0 {
		return nil
	}

	query := `
		INSERT INTO user_period_stats (stream_id, username, msg_count, period_start)
		VALUES (:stream_id, :username, :msg_count, :period_start)
		ON CONFLICT (stream_id, username)
		DO UPDATE SET 
		    msg_count = user_period_stats.msg_count + EXCLUDED.msg_count,
		    period_start = EXCLUDED.period_start; -- Обновляем время последнего сообщения
	`

	type userRow struct {
		StreamID    string    `db:"stream_id"`
		Username    string    `db:"username"`
		MsgCount    int       `db:"msg_count"`
		PeriodStart time.Time `db:"period_start"`
	}

	rows := make([]userRow, 0, len(stats.UserStats))
	for username, count := range stats.UserStats {
		rows = append(rows, userRow{
			StreamID:    stats.StreamID,
			Username:    username,
			MsgCount:    count,
			PeriodStart: stats.StartedAt,
		})
	}

	_, err := u.db.NamedExecContext(ctx, query, rows)
	return err
}
