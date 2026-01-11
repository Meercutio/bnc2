package game

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	RoundDuration time.Duration // 0 => таймер выключен
}

type Server struct {
	cfg Config

	mu      sync.Mutex
	matches map[string]*Match
}

func NewServer(cfg Config) *Server {
	return &Server{
		cfg:     cfg,
		matches: make(map[string]*Match),
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	// создать матч
	mux.HandleFunc("/api/match", s.handleCreateMatch)

	// websocket
	mux.HandleFunc("/ws", s.handleWS)

	// статика
	mux.Handle("/", http.FileServer(http.Dir("./web")))
}

func (s *Server) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	matchID := randID(10)

	m := NewMatch(matchID, s.cfg.RoundDuration)

	s.mu.Lock()
	s.matches[matchID] = m
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{
		"matchId": matchID,
	})
}

func (s *Server) getMatch(matchID string) (*Match, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.matches[matchID]
	return m, ok
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
