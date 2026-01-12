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
	cfg     Config
	matches *MatchService
}

func NewServer(cfg Config, matches *MatchService) *Server {
	return &Server{
		cfg:     cfg,
		matches: matches,
	}
}

// (опционально) если хочешь подменять storage в тестах/будущем:
func NewServerWithStore(cfg Config, matches *MatchService) *Server {
	return &Server{
		cfg:     cfg,
		matches: matches,
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

	_, err := s.matches.Create(r.Context(), matchID)
	if err != nil {
		http.Error(w, "failed to create match", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"matchId": matchID,
	})
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
