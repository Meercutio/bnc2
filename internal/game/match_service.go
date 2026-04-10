package game

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// MatchService отвечает за:
// - in-memory кэш матчей
// - восстановление матчей из persistent storage (Redis)
type MatchService struct {
	mu sync.Mutex
	in map[string]*Match

	cfg     Config
	persist MatchPersistence
	sink    GameResultSink
}

func NewMatchService(cfg Config, persist MatchPersistence) *MatchService {
	return &MatchService{
		in:      make(map[string]*Match),
		cfg:     cfg,
		persist: persist,
	}
}

func (s *MatchService) SetResultSink(sink GameResultSink) {
	s.mu.Lock()
	s.sink = sink
	s.mu.Unlock()
}

const persistTimeout = 2 * time.Second

func (s *MatchService) persistOnce(matchID string, snap MatchSnapshot) error {
	ctx, cancel := context.WithTimeout(context.Background(), persistTimeout)
	defer cancel()
	return s.persist.Save(ctx, matchID, snap)
}

func (s *MatchService) persistBestEffort(matchID string, snap MatchSnapshot) {
	if err := s.persistOnce(matchID, snap); err != nil {
		slog.Default().Warn("match persist failed", "matchId", matchID, "err", err)
	}
}

func (s *MatchService) CreateRanked(ctx context.Context, matchID string) (*Match, error) {
	return s.createWithRanked(ctx, matchID, true)
}

func (s *MatchService) Create(ctx context.Context, matchID string) (*Match, error) {
	return s.createWithRanked(ctx, matchID, false)
}

func (s *MatchService) createWithRanked(ctx context.Context, matchID string, ranked bool) (*Match, error) {
	m := NewMatch(matchID, s.cfg.RoundDuration)
	m.ranked = ranked

	// best-effort, async game-finished sink
	m.onGameFinished = func(e GameFinishedEvent) {
		s.mu.Lock()
		sink := s.sink
		s.mu.Unlock()
		if sink != nil {
			sink.HandleGameFinished(e)
		}
	}

	// hook: любое изменение матча будет сохранять snapshot
	m.onPersist = func(snap MatchSnapshot) {
		s.persistBestEffort(matchID, snap)
	}

	// первичное сохранение — важно, чтобы матч можно было восстановить после рестарта
	m.mu.Lock()
	snap := m.snapshotLocked()
	m.mu.Unlock()

	if err := s.persistOnce(matchID, snap); err != nil {
		return nil, fmt.Errorf("persist initial snapshot: %w", err)
	}

	return m, nil
}

func (s *MatchService) GetOrLoad(ctx context.Context, matchID string) (*Match, bool, error) {
	s.mu.Lock()
	m, ok := s.in[matchID]
	s.mu.Unlock()
	if ok {
		return m, true, nil
	}

	snap, found, err := s.persist.Load(ctx, matchID)
	if err != nil || !found {
		return nil, false, err
	}

	m = NewMatch(matchID, s.cfg.RoundDuration)
	m.onGameFinished = func(e GameFinishedEvent) {
		s.mu.Lock()
		sink := s.sink
		s.mu.Unlock()
		if sink != nil {
			sink.HandleGameFinished(e)
		}
	}

	// hook снова навешиваем (ВАЖНО: без request-context, иначе persist "отвалится" после завершения запроса)
	m.onPersist = func(snap MatchSnapshot) {
		s.persistBestEffort(matchID, snap)
	}

	m.mu.Lock()
	m.restoreLocked(snap)

	// reconcile: если матч в playing, но дедлайн уже прошёл — не зависаем в ожидании guess.
	if s.cfg.RoundDuration > 0 && m.phase == "playing" && m.roundActive && !m.deadline.IsZero() {
		now := time.Now()
		if !now.Before(m.deadline) {
			// дедлайн в прошлом — считаем пропуск для тех, кто не успел
			if !m.p1.guessSet {
				m.p1.missed = true
			}
			if !m.p2.guessSet {
				m.p2.missed = true
			}
			m.broadcastStateLocked()
			m.persistLocked()
			m.finalizeRoundLocked()
		} else {
			// дедлайн в будущем — поднимаем таймер заново
			m.roundToken++
			token := m.roundToken

			if m.roundTimer != nil {
				m.roundTimer.Stop()
			}
			d := time.Until(m.deadline)
			m.roundTimer = time.AfterFunc(d, func() {
				m.onRoundTimeout(token)
			})
		}
	}
	m.mu.Unlock()

	// If the match is already finished (e.g. restored after server restart) — emit again.
	// Downstream MUST be idempotent by matchId.
	m.mu.Lock()
	m.emitGameFinishedLocked()
	m.mu.Unlock()

	s.mu.Lock()
	s.in[matchID] = m
	s.mu.Unlock()

	return m, true, nil
}

func (s *MatchService) ReleaseIfIdle(matchID string, m *Match) {
	if m == nil || m.ConnectedCount() > 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if current, ok := s.in[matchID]; ok && current == m && m.ConnectedCount() == 0 {
		delete(s.in, matchID)
	}
}
