package game

import (
	"context"
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
}

func NewMatchService(cfg Config, persist MatchPersistence) *MatchService {
	return &MatchService{
		in:      make(map[string]*Match),
		cfg:     cfg,
		persist: persist,
	}
}

func (s *MatchService) Create(ctx context.Context, matchID string) (*Match, error) {
	m := NewMatch(matchID, s.cfg.RoundDuration)

	// hook: любое изменение матча будет сохранять snapshot
	m.onPersist = func(snap MatchSnapshot) {
		_ = s.persist.Save(ctx, matchID, snap) // MVP: без логирования
	}

	// первичное сохранение
	m.mu.Lock()
	snap := m.snapshotLocked()
	m.mu.Unlock()
	_ = s.persist.Save(ctx, matchID, snap)

	s.mu.Lock()
	s.in[matchID] = m
	s.mu.Unlock()

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
	m.mu.Lock()
	m.restoreLocked(snap)
	m.mu.Unlock()

	// hook снова навешиваем
	m.onPersist = func(snap MatchSnapshot) {
		_ = s.persist.Save(ctx, matchID, snap)
	}

	// если матч в playing и дедлайн ещё не прошёл — поднимаем таймер заново
	m.mu.Lock()
	if s.cfg.RoundDuration > 0 && m.phase == "playing" && m.roundActive && !m.deadline.IsZero() && time.Now().Before(m.deadline) {
		// новый token, чтобы старые таймеры (до рестарта) не влияли
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
	m.mu.Unlock()

	s.mu.Lock()
	s.in[matchID] = m
	s.mu.Unlock()

	return m, true, nil
}
