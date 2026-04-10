package game

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"example.com/bc-mvp/internal/auth"
)

type Config struct {
	RoundDuration time.Duration // 0 => таймер выключен
}

type Server struct {
	cfg     Config
	matches *MatchService
	auth    TokenVerifier
	mm      *Matchmaker
}

type TokenVerifier interface {
	Verify(token string) (*auth.Claims, error)
}

func NewServer(cfg Config, matches *MatchService, verifier TokenVerifier, mm *Matchmaker) *Server {
	return &Server{
		cfg:     cfg,
		matches: matches,
		auth:    verifier,
		mm:      mm,
	}
}

// (опционально) если хочешь подменять storage в тестах/будущем:
//func NewServerWithStore(cfg Config, matches *MatchService) *Server {
//	return &Server{
//		cfg:     cfg,
//		matches: matches,
//	}
//}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/match", s.handleCreateMatch)

	mux.HandleFunc("/api/matchmaking/join", s.handleMatchmakingJoin)
	mux.HandleFunc("/api/matchmaking/stats", s.handleMatchmakingStats)

	mux.HandleFunc("/ws/", s.handleWS)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing matchId: use /ws/{matchId}", http.StatusBadRequest)
	})
}

func (s *Server) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if _, err := s.verifyRequestAuth(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
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

func (s *Server) handleMatchmakingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.mm == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]any{"error": "matchmaking disabled"})
		return
	}
	waiting, forDur := s.mm.Stats()
	writeJSON(w, http.StatusOK, map[string]any{
		"waiting":      waiting,
		"waitingForMs": forDur.Milliseconds(),
	})
}

func (s *Server) handleMatchmakingJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.mm == nil {
		writeJSON(w, http.StatusNotImplemented, map[string]any{"error": "matchmaking disabled"})
		return
	}

	claims, err := s.verifyRequestAuth(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	// Long-poll: keep request open until matched or client aborts.
	matchID, jerr := s.mm.Join(r.Context(), claims.UserID, claims.DisplayName, func() (string, error) {
		id := randID(10)
		_, err := s.matches.CreateRanked(r.Context(), id)
		if err != nil {
			return "", err
		}
		return id, nil
	})

	if jerr != nil {
		if errorsIsContextDone(jerr) {
			// если клиент отменил запрос — часто ответа не увидит, но пусть будет корректно
			writeJSON(w, http.StatusAccepted, map[string]any{"status": "canceled"})
			return
		}
		if jerr == ErrAlreadyWaiting {
			writeJSON(w, http.StatusConflict, map[string]any{"error": "already_waiting"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "matchmaking_error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "matched",
		"matchId": matchID,
	})
}

func (s *Server) verifyRequestAuth(r *http.Request) (*auth.Claims, error) {
	tok := strings.TrimSpace(auth.TokenFromRequest(r))
	if tok == "" {
		return nil, errors.New("missing auth token")
	}
	return s.auth.Verify(tok)
}

func errorsIsContextDone(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
