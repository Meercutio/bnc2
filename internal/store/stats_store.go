package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerStats struct {
	UserID    string
	Wins      int
	Losses    int
	Draws     int
	UpdatedAt time.Time
}

type StatsStore struct {
	db *pgxpool.Pool
}

func NewStatsStore(db *pgxpool.Pool) *StatsStore {
	return &StatsStore{db: db}
}

func (s *StatsStore) InitForUser(ctx context.Context, userID string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO player_stats (user_id, wins, losses, draws)
		VALUES ($1, 0, 0, 0)
		ON CONFLICT (user_id) DO NOTHING
	`, userID)
	return err
}

func (s *StatsStore) Get(ctx context.Context, userID string) (PlayerStats, error) {
	var st PlayerStats
	err := s.db.QueryRow(ctx, `
		SELECT user_id, wins, losses, draws, updated_at
		FROM player_stats
		WHERE user_id=$1
	`, userID).Scan(&st.UserID, &st.Wins, &st.Losses, &st.Draws, &st.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		// если вдруг статистики нет — это не фатально, можно считать нулями
		return PlayerStats{UserID: userID}, nil
	}
	if err != nil {
		return PlayerStats{}, err
	}
	return st, nil
}
