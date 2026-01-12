package game

import "sync"

// MatchStore — контракт хранения матчей.
// Сейчас будет in-memory, позже можно сделать Redis/Postgres реализацию.
type MatchStore interface {
	Create(matchID string, m *Match)
	Get(matchID string) (*Match, bool)
	Delete(matchID string)
}

type InMemoryMatchStore struct {
	mu sync.Mutex
	m  map[string]*Match
}

func NewInMemoryMatchStore() *InMemoryMatchStore {
	return &InMemoryMatchStore{
		m: make(map[string]*Match),
	}
}

func (s *InMemoryMatchStore) Create(matchID string, m *Match) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[matchID] = m
}

func (s *InMemoryMatchStore) Get(matchID string) (*Match, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.m[matchID]
	return m, ok
}

func (s *InMemoryMatchStore) Delete(matchID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, matchID)
}
