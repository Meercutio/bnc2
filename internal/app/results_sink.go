package app

import (
	"context"
	"log/slog"
	"time"

	"example.com/bc-mvp/internal/game"
	"example.com/bc-mvp/internal/store"
)

// PGGameResultSink writes finished match results into Postgres (W/L/D + rating).
// It is best-effort; processing is idempotent by match_id.
//
// NOTE: For MVP we process in-process. If you'll add Kafka later,
// this is the exact place to publish events instead.
type PGGameResultSink struct {
	Stats *store.StatsStore
	Log   *slog.Logger
}

func (s *PGGameResultSink) HandleGameFinished(e game.GameFinishedEvent) {
	if s == nil || s.Stats == nil {
		return
	}
	if !e.Ranked {
		// Friendly matches do not affect stats/rating.
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stats.ApplyGameResult(ctx, e.ResultID, e.MatchID, e.GameNo, e.P1ID, e.P2ID, e.Winner); err != nil {
		if s.Log != nil {
			s.Log.Warn("apply game result failed", "matchId", e.MatchID, "err", err)
		}
	}
}
