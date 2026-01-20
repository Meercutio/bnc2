//go:build integration

package game

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func newRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{Addr: addr})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	require.NoError(t, rdb.Ping(ctx).Err(), "redis is not reachable")
	return rdb
}

func TestRedisPersistence_CreateSaveLoad(t *testing.T) {
	ctx := context.Background()
	rdb := newRedisClient(t)

	// Чистим Redis, чтобы тест был детерминированный
	require.NoError(t, rdb.FlushDB(ctx).Err())

	ttl := 1 * time.Hour
	persist := NewRedisMatchStore(rdb, ttl)

	cfg := Config{RoundDuration: 0}
	svc1 := NewMatchService(cfg, persist)

	matchID := "m_test_1"

	// 1) Создали матч и сохранили snapshot
	_, err := svc1.Create(ctx, matchID)
	require.NoError(t, err)

	// 2) В памяти "поиграли": 2 игрока, секреты, ход
	m, ok, err := svc1.GetOrLoad(ctx, matchID)
	require.NoError(t, err)
	require.True(t, ok)

	// attach игроков
	_, code, _ := m.Attach("u1", "Alice", newTestConn())
	require.Empty(t, code)
	_, code, _ = m.Attach("u2", "Bob", newTestConn())
	require.Empty(t, code)

	require.NoError(t, m.SetSecret(P1, "0011"))
	require.NoError(t, m.SetSecret(P2, "0101"))

	require.NoError(t, m.SubmitGuess(P1, "0101"))
	require.NoError(t, m.SubmitGuess(P2, "0011"))

	// 3) Симулируем рестарт: новый MatchService с пустым in-memory
	svc2 := NewMatchService(cfg, persist)
	m2, ok, err := svc2.GetOrLoad(ctx, matchID)
	require.NoError(t, err)
	require.True(t, ok)

	// 4) Проверяем, что состояние восстановилось
	m2.mu.Lock()
	defer m2.mu.Unlock()

	require.Equal(t, "finished", m2.phase)
	require.Equal(t, 1, m2.round)
	require.Equal(t, "draw", m2.winner)
	require.Len(t, m2.history, 1)
}

func TestRedisPersistence_RestoreActiveRound_TimerOff(t *testing.T) {
	ctx := context.Background()
	rdb := newRedisClient(t)
	require.NoError(t, rdb.FlushDB(ctx).Err())

	ttl := 1 * time.Hour
	persist := NewRedisMatchStore(rdb, ttl)

	// таймер выключен
	cfg := Config{RoundDuration: 0}
	svc := NewMatchService(cfg, persist)

	matchID := "m_test_2"
	m, err := svc.Create(ctx, matchID)
	require.NoError(t, err)

	_, code, _ := m.Attach("u1", "Alice", newTestConn())
	require.Empty(t, code)
	_, code, _ = m.Attach("u2", "Bob", newTestConn())
	require.Empty(t, code)

	require.NoError(t, m.SetSecret(P1, "1111"))
	require.NoError(t, m.SetSecret(P2, "2222"))

	// Теперь матч должен быть playing/roundActive, и это должно сохраниться
	// Рестарт:
	svc2 := NewMatchService(cfg, persist)
	m2, ok, err := svc2.GetOrLoad(ctx, matchID)
	require.NoError(t, err)
	require.True(t, ok)

	m2.mu.Lock()
	defer m2.mu.Unlock()

	require.Equal(t, "playing", m2.phase)
	require.Equal(t, 1, m2.round)
	require.True(t, m2.roundActive)
}
