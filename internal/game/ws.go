package game

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

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
//
// JWT больше не передаём в query-string.
// Поддерживаются 2 варианта:
//  1. Authorization: Bearer <jwt> (для клиентов, которые умеют ставить headers)
//  2. Первое WS-сообщение: {"type":"auth","payload":{"token":"..."}} (для browser WebSocket)
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	matchID, ok := matchIDFromWSPath(r.URL.Path)
	if !ok {
		http.Error(w, "missing or invalid matchId: use /ws/{matchId}", http.StatusBadRequest)
		return
	}

	// получаем матч (in-memory или из Redis) ДО upgrade, чтобы быстрее отсеять 404
	m, ok, err := s.matches.GetOrLoad(r.Context(), matchID)
	if err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "match not found", http.StatusNotFound)
		return
	}

	// Вариант 1: token из headers
	playerID, displayName, err := s.authFromRequest(r)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Если токена не было в headers — ожидаем auth-сообщение как первое.
	if playerID == "" {
		pid, name, aerr := s.authOverWS(ws)
		if aerr != nil {
			_ = ws.WriteJSON(Envelope{Type: "error", Payload: mustJSON(ErrorPayload{Code: "unauthorized", Message: aerr.Error()})})
			_ = ws.Close()
			return
		}
		playerID = pid
		displayName = name
	}

	cc := &ClientConn{
		ws:   ws,
		send: make(chan []byte, 64),
	}

	slot, errCode, errMsg := m.Attach(playerID, displayName, cc)
	if errCode != "" {
		_ = ws.WriteJSON(Envelope{
			Type:    "error",
			Payload: mustJSON(ErrorPayload{Code: errCode, Message: errMsg}),
		})
		cc.Close()
		return
	}

	// writer loop (теперь уже после успешной авторизации)
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
		case "auth":
			m.SendErrorTo(slot, "bad_input", "already authorized")
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

func (s *Server) authFromRequest(r *http.Request) (userID string, displayName string, err error) {
	// Authorization: Bearer <token>
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		tok := strings.TrimPrefix(h, "Bearer ")
		claims, err := s.auth.Verify(tok)
		if err != nil {
			return "", "", err
		}
		return claims.UserID, claims.DisplayName, nil
	}

	// Sec-WebSocket-Protocol: <token> (используют некоторые клиенты)
	// В браузере `new WebSocket(url, [token])` попадает сюда.
	if p := strings.TrimSpace(r.Header.Get("Sec-WebSocket-Protocol")); p != "" {
		// может быть список протоколов через запятую
		parts := strings.Split(p, ",")
		for _, part := range parts {
			tok := strings.TrimSpace(part)
			if tok == "" {
				continue
			}
			claims, err := s.auth.Verify(tok)
			if err == nil {
				return claims.UserID, claims.DisplayName, nil
			}
		}
	}

	// Не ошибка: просто придётся авторизоваться через первое WS-сообщение.
	return "", "", nil
}

type authPayload struct {
	Token string `json:"token"`
}

func (s *Server) authOverWS(ws *websocket.Conn) (userID string, displayName string, err error) {
	_ = ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer ws.SetReadDeadline(time.Time{})

	_, data, err := ws.ReadMessage()
	if err != nil {
		return "", "", err
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return "", "", err
	}
	if env.Type != "auth" {
		return "", "", errors.New("missing auth message")
	}
	var p authPayload
	if err := json.Unmarshal(env.Payload, &p); err != nil {
		return "", "", err
	}
	if strings.TrimSpace(p.Token) == "" {
		return "", "", errors.New("missing token")
	}
	claims, err := s.auth.Verify(strings.TrimSpace(p.Token))
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.DisplayName, nil
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

// matchIDFromWSPath extracts matchId from /ws/{matchId}.
//
// net/http ServeMux matches "/ws/" as a subtree, so we must validate and reject paths like "/ws/abc/def".
func matchIDFromWSPath(path string) (string, bool) {
	const prefix = "/ws/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	id := strings.TrimPrefix(path, prefix)
	if id == "" {
		return "", false
	}
	// no extra path segments
	if strings.Contains(id, "/") {
		return "", false
	}
	// simple validation: match IDs are generated as [a-z0-9]+ by randID
	if len(id) > 64 {
		return "", false
	}
	for _, ch := range id {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			continue
		}
		return "", false
	}
	return id, true
}
