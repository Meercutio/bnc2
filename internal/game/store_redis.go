package game

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// MatchPersistence — абстракция "положить/достать snapshot".
// Реализуем Redis-ом. Позже можно добавить Postgres.
type MatchPersistence interface {
	Save(ctx context.Context, matchID string, snap MatchSnapshot) error
	Load(ctx context.Context, matchID string) (MatchSnapshot, bool, error)
}

type RedisMatchStore struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisMatchStore(rdb *redis.Client, ttl time.Duration) *RedisMatchStore {
	return &RedisMatchStore{rdb: rdb, ttl: ttl}
}

func (s *RedisMatchStore) key(matchID string) string {
	return fmt.Sprintf("match:%s:snapshot", matchID)
}

func (s *RedisMatchStore) Save(ctx context.Context, matchID string, snap MatchSnapshot) error {
	b, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, s.key(matchID), b, s.ttl).Err()
}

func (s *RedisMatchStore) Load(ctx context.Context, matchID string) (MatchSnapshot, bool, error) {
	val, err := s.rdb.Get(ctx, s.key(matchID)).Bytes()
	if err == redis.Nil {
		return MatchSnapshot{}, false, nil
	}
	if err != nil {
		return MatchSnapshot{}, false, err
	}

	var snap MatchSnapshot
	if err := json.Unmarshal(val, &snap); err != nil {
		return MatchSnapshot{}, false, err
	}
	return snap, true, nil
}
