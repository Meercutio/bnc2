package game

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"example.com/bc-mvp/internal/auth"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // MVP
}

type ClientConn struct {
	ws   *websocket.Conn
	send chan []byte

	closeOnce sync.Once
}

func (c *ClientConn) Close() {
	c.closeOnce.Do(func() {
		close(c.send)
		_ = c.ws.Close()
	})
}

// handleWS — WebSocket вход в матч
// Требует JWT: /ws?matchId=xxx&token=yyy
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	matchID := r.URL.Query().Get("matchId")
	token := r.URL.Query().Get("token")

	if matchID == "" || token == "" {
		http.Error(w, "missing matchId or token", http.StatusBadRequest)
		return
	}

	// 🔐 Проверяем JWT
	claims, err := auth.Verify(jwtSecret(), token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	playerID := claims.UserID

	// получаем матч (in-memory или из Redis)
	m, ok, err := s.matches.GetOrLoad(r.Context(), matchID)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "match not found", http.StatusNotFound)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	cc := &ClientConn{
		ws:   ws,
		send: make(chan []byte, 64),
	}

	slot, errCode, errMsg := m.Attach(playerID, cc)
	if errCode != "" {
		_ = ws.WriteJSON(Envelope{
			Type:    "error",
			Payload: mustJSON(ErrorPayload{Code: errCode, Message: errMsg}),
		})
		cc.Close()
		return
	}

	// writer loop
	go func() {
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-cc.send:
				if !ok {
					return
				}
				_ = ws.WriteMessage(websocket.TextMessage, msg)
			case <-ticker.C:
				_ = ws.WriteMessage(websocket.PingMessage, []byte{})
			}
		}
	}()

	// initial state
	m.SendStateTo(slot)
	m.BroadcastState()

	// reader loop
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var env Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			m.SendErrorTo(slot, "bad_json", "invalid json")
			continue
		}

		switch env.Type {
		case "set_secret":
			var p SetSecretPayload
			if err := json.Unmarshal(env.Payload, &p); err != nil {
				m.SendErrorTo(slot, "bad_input", "invalid payload")
				continue
			}
			if err := m.SetSecret(slot, p.Secret); err != nil {
				m.SendErrorTo(slot, "bad_input", err.Error())
			}

		case "submit_guess":
			var p SubmitGuessPayload
			if err := json.Unmarshal(env.Payload, &p); err != nil {
				m.SendErrorTo(slot, "bad_input", "invalid payload")
				continue
			}
			if err := m.SubmitGuess(slot, p.Guess); err != nil {
				m.SendErrorTo(slot, "bad_input", err.Error())
			}

		case "rematch_request":
			if err := m.RequestRematch(slot); err != nil {
				m.SendErrorTo(slot, "bad_input", err.Error())
			}

		default:
			m.SendErrorTo(slot, "unknown_type", "unknown message type")
		}
	}

	// disconnect
	m.Detach(slot)
	cc.Close()
	m.BroadcastState()
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
