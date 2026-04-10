package store

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerStats struct {
	UserID    string
	Wins      int
	Losses    int
	Draws     int
	Games     int
	Rating    int
	UpdatedAt time.Time
}

type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	Rating      int    `json:"rating"`
	Games       int    `json:"games"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	Draws       int    `json:"draws"`
}

type MyRating struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	Rating      int    `json:"rating"`
	Games       int    `json:"games"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	Draws       int    `json:"draws"`
	Rank        int    `json:"rank"`
}

type StatsStore struct {
	db *pgxpool.Pool
}

func NewStatsStore(db *pgxpool.Pool) *StatsStore {
	return &StatsStore{db: db}
}

func (s *StatsStore) InitForUser(ctx context.Context, userID string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO player_stats (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`, userID)
	return err
}

func (s *StatsStore) Get(ctx context.Context, userID string) (PlayerStats, error) {
	var st PlayerStats
	err := s.db.QueryRow(ctx, `
		SELECT user_id, wins, losses, draws, games, rating, updated_at
		FROM player_stats
		WHERE user_id=$1
	`, userID).Scan(&st.UserID, &st.Wins, &st.Losses, &st.Draws, &st.Games, &st.Rating, &st.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		if ensureErr := s.InitForUser(ctx, userID); ensureErr != nil {
			return PlayerStats{}, ensureErr
		}
		return PlayerStats{UserID: userID, Rating: 1000}, nil
	}
	if err != nil {
		return PlayerStats{}, err
	}
	return st, nil
}

func (s *StatsStore) Leaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	rows, err := s.db.Query(ctx, `
		SELECT
			RANK() OVER (ORDER BY ps.rating DESC, ps.games DESC, ps.wins DESC, u.display_name ASC) AS rnk,
			u.id, u.display_name, ps.rating, ps.games, ps.wins, ps.losses, ps.draws
		FROM player_stats ps
		JOIN users u ON u.id = ps.user_id
		ORDER BY ps.rating DESC, ps.games DESC, ps.wins DESC, u.display_name ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]LeaderboardEntry, 0, limit)
	for rows.Next() {
		var e LeaderboardEntry
		err := rows.Scan(&e.Rank, &e.UserID, &e.DisplayName, &e.Rating, &e.Games, &e.Wins, &e.Losses, &e.Draws)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *StatsStore) MyRating(ctx context.Context, userID string) (MyRating, error) {
	var out MyRating
	err := s.db.QueryRow(ctx, `
		WITH ranked AS (
			SELECT
				ps.user_id,
				RANK() OVER (ORDER BY ps.rating DESC, ps.games DESC, ps.wins DESC, u.display_name ASC) AS rnk,
				u.display_name,
				ps.rating,
				ps.games,
				ps.wins,
				ps.losses,
				ps.draws
			FROM player_stats ps
			JOIN users u ON u.id = ps.user_id
		)
		SELECT user_id, display_name, rating, games, wins, losses, draws, rnk
		FROM ranked
		WHERE user_id = $1
	`, userID).Scan(&out.UserID, &out.DisplayName, &out.Rating, &out.Games, &out.Wins, &out.Losses, &out.Draws, &out.Rank)
	if errors.Is(err, pgx.ErrNoRows) {
		if ensureErr := s.InitForUser(ctx, userID); ensureErr != nil {
			return MyRating{}, ensureErr
		}
		return MyRating{UserID: userID, Rating: 1000, Rank: 0}, nil
	}
	return out, err
}

// ApplyGameResult updates W/L/D + Elo rating for both players.
// It is idempotent per matchID (via game_results.match_id unique).
func (s *StatsStore) ApplyGameResult(ctx context.Context, resultID, matchID string, gameNo int, p1ID, p2ID, winner string) error {
	// Only registered users (UUID). If not UUID — skip silently (MVP).
	p1u, err := uuid.Parse(p1ID)
	if err != nil {
		return nil
	}
	p2u, err := uuid.Parse(p2ID)
	if err != nil {
		return nil
	}
	if matchID == "" || resultID == "" {
		return fmt.Errorf("matchID/resultID is empty")
	}
	if gameNo <= 0 {
		gameNo = 1
	}

	// outcome scores
	var s1, s2 float64
	switch winner {
	case "p1":
		s1, s2 = 1, 0
	case "p2":
		s1, s2 = 0, 1
	case "draw":
		s1, s2 = 0.5, 0.5
	default:
		return fmt.Errorf("invalid winner: %s", winner)
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// ensure stats rows exist
	_, _ = tx.Exec(ctx, `
		INSERT INTO player_stats (user_id, wins, losses, draws, games, rating)
		VALUES ($1,0,0,0,0,1000)
		ON CONFLICT (user_id) DO NOTHING
	`, p1u)
	_, _ = tx.Exec(ctx, `
		INSERT INTO player_stats (user_id, wins, losses, draws, games, rating)
		VALUES ($1,0,0,0,0,1000)
		ON CONFLICT (user_id) DO NOTHING
	`, p2u)

	// lock both rows
	var r1, r2 int
	var g1, g2 int
	var w1, l1, d1 int
	var w2, l2, d2 int

	if err := tx.QueryRow(ctx, `
		SELECT rating, games, wins, losses, draws
		FROM player_stats
		WHERE user_id=$1
		FOR UPDATE
	`, p1u).Scan(&r1, &g1, &w1, &l1, &d1); err != nil {
		return err
	}
	if err := tx.QueryRow(ctx, `
		SELECT rating, games, wins, losses, draws
		FROM player_stats
		WHERE user_id=$1
		FOR UPDATE
	`, p2u).Scan(&r2, &g2, &w2, &l2, &d2); err != nil {
		return err
	}

	k1 := kFactor(r1, g1)
	k2 := kFactor(r2, g2)

	e1 := expectedScore(r1, r2)
	e2 := expectedScore(r2, r1)

	nr1 := int(math.Round(float64(r1) + k1*(s1-e1)))
	nr2 := int(math.Round(float64(r2) + k2*(s2-e2)))

	// idempotency gate: insert game_results first.
	ct, err := tx.Exec(ctx, `
		INSERT INTO game_results(
			result_id, match_id, game_no,
			p1_id, p2_id, winner,
			p1_rating_before, p1_rating_after,
			p2_rating_before, p2_rating_after
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (result_id) DO NOTHING
	`, resultID, matchID, gameNo, p1u, p2u, winner, r1, nr1, r2, nr2)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		// already processed
		return tx.Commit(ctx)
	}

	// apply W/L/D
	if winner == "p1" {
		w1++
		l2++
	} else if winner == "p2" {
		w2++
		l1++
	} else {
		d1++
		d2++
	}
	g1++
	g2++

	_, err = tx.Exec(ctx, `
		UPDATE player_stats
		SET wins=$2, losses=$3, draws=$4, games=$5, rating=$6, updated_at=now()
		WHERE user_id=$1
	`, p1u, w1, l1, d1, g1, nr1)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		UPDATE player_stats
		SET wins=$2, losses=$3, draws=$4, games=$5, rating=$6, updated_at=now()
		WHERE user_id=$1
	`, p2u, w2, l2, d2, g2, nr2)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func expectedScore(rA, rB int) float64 {
	return 1.0 / (1.0 + math.Pow(10.0, float64(rB-rA)/400.0))
}

func kFactor(rating, games int) float64 {
	// MVP: классическая схема, чтобы новички быстрее "находили" свой уровень.
	// - первые 30 игр: K=32
	// - дальше: K=16
	// - для высоких рейтингов (>=2000): K=12
	if games < 30 {
		return 32
	}
	if rating >= 2000 {
		return 12
	}
	return 16
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
