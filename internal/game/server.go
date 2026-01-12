package game

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"
)

type Config struct {
	RoundDuration time.Duration // 0 => таймер выключен
}

type Server struct {
	cfg   Config
	store MatchStore
}

func NewServer(cfg Config) *Server {
	return &Server{
		cfg:   cfg,
		store: NewInMemoryMatchStore(), // MVP default
	}
}

// (опционально) если хочешь подменять storage в тестах/будущем:
func NewServerWithStore(cfg Config, store MatchStore) *Server {
	return &Server{
		cfg:   cfg,
		store: store,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/match", s.handleCreateMatch)
	mux.HandleFunc("/ws", s.handleWS)
}

func (s *Server) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	matchID := randID(10)
	m := NewMatch(matchID, s.cfg.RoundDuration)

	s.store.Create(matchID, m)

	writeJSON(w, http.StatusOK, map[string]string{
		"matchId": matchID,
	})
}

func (s *Server) getMatch(matchID string) (*Match, bool) {
	return s.store.Get(matchID)
}

func randID(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
