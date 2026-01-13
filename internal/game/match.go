package game

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type Slot string

const (
	P1 Slot = "p1"
	P2 Slot = "p2"
)

type Match struct {
	id string
	mu sync.Mutex

	phase string // waiting_players|waiting_secrets|playing|finished

	round       int
	deadline    time.Time
	roundActive bool
	roundTimer  *time.Timer
	roundToken  int64
	roundDur    time.Duration
	winner      string // p1|p2|draw|""

	p1 *Player
	p2 *Player

	history      []RoundHistoryItem
	seriesP1Wins int
	seriesP2Wins int
	seriesDraws  int
	onPersist    func(MatchSnapshot)
}

type Player struct {
	id   string
	conn *ClientConn

	connected        bool
	rematchRequested bool

	secret    string
	secretSet bool

	guess    string
	guessSet bool
	missed   bool
}

func NewMatch(id string, roundDur time.Duration) *Match {
	return &Match{
		id:       id,
		phase:    "waiting_players",
		roundDur: roundDur,
		p1:       &Player{},
		p2:       &Player{},
	}
}

func (m *Match) Attach(playerID string, cc *ClientConn) (Slot, string, string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// reconnect?
	if m.p1.id == playerID && m.p1.id != "" {
		m.p1.conn = cc
		m.p1.connected = true
		return P1, "", ""
	}
	if m.p2.id == playerID && m.p2.id != "" {
		m.p2.conn = cc
		m.p2.connected = true
		return P2, "", ""
	}

	// new join
	if m.p1.id == "" {
		m.p1.id = playerID
		m.p1.conn = cc
		m.p1.connected = true
		m.updatePhaseLocked()
		return P1, "", ""
	}
	if m.p2.id == "" && m.p1.id != playerID {
		m.p2.id = playerID
		m.p2.conn = cc
		m.p2.connected = true
		m.updatePhaseLocked()
		return P2, "", ""
	}

	return "", "match_full", "match already has two players"
}

func (m *Match) Detach(slot Slot) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.playerLocked(slot)
	p.connected = false
	p.conn = nil
	m.updatePhaseLocked()
}

func (m *Match) SetSecret(slot Slot, secret string) error {
	if !valid4Digits(secret) {
		return errors.New("secret must be exactly 4 digits (0-9)")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.phase == "finished" {
		return errors.New("game already finished")
	}

	p := m.playerLocked(slot)
	p.secret = secret
	p.secretSet = true

	m.updatePhaseLocked()

	// ✅ ВАЖНО: стартуем раунд по факту готовности секретов,
	// а не по значению phase, потому что phase уже могла стать "playing".
	if m.p1.secretSet && m.p2.secretSet && !m.roundActive && m.round == 0 && m.phase != "finished" {
		m.phase = "playing"
		m.startRoundLocked()
	}

	m.broadcastStateLocked()
	m.persistLocked()
	return nil
}

func (m *Match) SubmitGuess(slot Slot, guess string) error {
	if !valid4Digits(guess) {
		return errors.New("guess must be exactly 4 digits (0-9)")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.phase != "playing" {
		return errors.New("game is not in playing phase")
	}
	if !m.roundActive {
		return errors.New("round is not active")
	}

	p := m.playerLocked(slot)
	if p.guessSet || p.missed {
		return errors.New("guess already submitted (or missed)")
	}

	p.guess = guess
	p.guessSet = true

	m.broadcastStateLocked()

	// если оба ввели — закрываем раунд
	if (m.p1.guessSet || m.p1.missed) && (m.p2.guessSet || m.p2.missed) {
		m.finalizeRoundLocked()
	}
	// если раунд ещё не закрыт — сохраняем частичное состояние (что один игрок уже ввёл)
	m.persistLocked()
	return nil
}

func (m *Match) RequestRematch(slot Slot) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.phase != "finished" {
		return errors.New("rematch available only after game finished")
	}

	p := m.playerLocked(slot)
	p.rematchRequested = true

	// сообщаем состояние рематча
	m.broadcastLocked(Envelope{
		Type: "rematch_status",
		Payload: mustJSON(map[string]any{
			"p1": m.p1.rematchRequested,
			"p2": m.p2.rematchRequested,
		}),
	})

	// если оба согласились — стартуем новую игру
	if m.p1.rematchRequested && m.p2.rematchRequested {
		m.startRematchLocked()
	}

	m.persistLocked()
	return nil
}

func (m *Match) startRematchLocked() {
	// сбрасываем флаги рематча
	m.p1.rematchRequested = false
	m.p2.rematchRequested = false

	// останавливаем таймер на всякий случай
	if m.roundTimer != nil {
		m.roundTimer.Stop()
	}

	// сбрасываем состояние матча, но оставляем игроков и соединения
	m.phase = "waiting_secrets"
	m.winner = ""
	m.round = 0
	m.roundActive = false
	m.deadline = time.Time{}
	m.history = nil

	m.p1.secret = ""
	m.p2.secret = ""
	m.p1.secretSet = false
	m.p2.secretSet = false

	m.p1.guess = ""
	m.p2.guess = ""
	m.p1.guessSet = false
	m.p2.guessSet = false
	m.p1.missed = false
	m.p2.missed = false

	// уведомляем фронт
	m.broadcastLocked(Envelope{
		Type: "rematch_started",
		Payload: mustJSON(map[string]any{
			"series": map[string]int{
				"p1Wins": m.seriesP1Wins,
				"p2Wins": m.seriesP2Wins,
				"draws":  m.seriesDraws,
			},
		}),
	})

	m.broadcastStateLocked()
	m.persistLocked()
}

func (m *Match) SendErrorTo(slot Slot, code, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.playerLocked(slot)
	if p.conn == nil {
		return
	}
	m.sendLocked(p.conn, Envelope{
		Type:    "error",
		Payload: mustJSON(ErrorPayload{Code: code, Message: message}),
	})
}

func (m *Match) SendStateTo(slot Slot) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := m.playerLocked(slot)
	if p.conn == nil {
		return
	}
	state := m.buildStateLocked(slot)
	m.sendLocked(p.conn, Envelope{Type: "state", Payload: mustJSON(state)})
}

func (m *Match) BroadcastState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.broadcastStateLocked()
}

func (m *Match) broadcastStateLocked() {
	// персонализируем "you" (p1/p2)
	if m.p1.conn != nil {
		state := m.buildStateLocked(P1)
		m.sendLocked(m.p1.conn, Envelope{Type: "state", Payload: mustJSON(state)})
	}
	if m.p2.conn != nil {
		state := m.buildStateLocked(P2)
		m.sendLocked(m.p2.conn, Envelope{Type: "state", Payload: mustJSON(state)})
	}
}

func (m *Match) updatePhaseLocked() {
	if m.phase == "finished" {
		return
	}
	if m.p1.id == "" || m.p2.id == "" || !(m.p1.connected && m.p2.connected) {
		m.phase = "waiting_players"
		return
	}
	if !m.p1.secretSet || !m.p2.secretSet {
		m.phase = "waiting_secrets"
		return
	}
	if m.phase == "waiting_players" || m.phase == "waiting_secrets" {
		m.phase = "playing"
	}
}

func (m *Match) startRoundLocked() {
	m.round++
	m.roundActive = true

	// reset per-round
	m.p1.guessSet, m.p2.guessSet = false, false
	m.p1.missed, m.p2.missed = false, false
	m.p1.guess, m.p2.guess = "", ""

	// deadline/timer (итерация 2)
	if m.roundDur > 0 {
		m.deadline = time.Now().Add(m.roundDur)
		m.roundToken++
		token := m.roundToken

		if m.roundTimer != nil {
			m.roundTimer.Stop()
		}
		m.roundTimer = time.AfterFunc(m.roundDur, func() {
			m.onRoundTimeout(token)
		})
	} else {
		m.deadline = time.Time{}
	}

	// событие round_started
	payload := RoundStartedPayload{
		Round:      m.round,
		DeadlineMs: toMs(m.deadline),
	}
	m.broadcastLocked(Envelope{Type: "round_started", Payload: mustJSON(payload)})
	m.broadcastStateLocked()
}

func (m *Match) onRoundTimeout(token int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.phase != "playing" || !m.roundActive {
		return
	}
	if token != m.roundToken {
		return // старый таймер
	}

	// вариант A: у кого нет guess — пропуск
	if !m.p1.guessSet {
		m.p1.missed = true
	}
	if !m.p2.guessSet {
		m.p2.missed = true
	}

	m.broadcastStateLocked()
	m.persistLocked()
	m.finalizeRoundLocked()
}

func (m *Match) finalizeRoundLocked() {
	if !m.roundActive {
		return
	}
	m.roundActive = false

	// stop timer
	if m.roundTimer != nil {
		m.roundTimer.Stop()
	}

	// считаем результат для каждого: guess против секрета соперника
	a1 := m.attemptLocked(P1)
	a2 := m.attemptLocked(P2)

	item := RoundHistoryItem{
		Round: m.round,
		P1:    a1,
		P2:    a2,
	}
	m.history = append(m.history, item)

	// победа/ничья
	p1win := a1.Guess != nil && a1.Bulls == 4
	p2win := a2.Guess != nil && a2.Bulls == 4

	switch {
	case p1win && p2win:
		m.winner = "draw"
		m.phase = "finished"
	case p1win:
		m.winner = "p1"
		m.phase = "finished"
	case p2win:
		m.winner = "p2"
		m.phase = "finished"
	}

	if m.phase == "finished" {
		switch m.winner {
		case "p1":
			m.seriesP1Wins++
		case "p2":
			m.seriesP2Wins++
		case "draw":
			m.seriesDraws++
		}
	}

	// событие round_result
	m.broadcastLocked(Envelope{Type: "round_result", Payload: mustJSON(item)})
	m.broadcastStateLocked()

	if m.phase == "finished" {
		m.broadcastLocked(Envelope{
			Type: "series_score",
			Payload: mustJSON(map[string]any{
				"series": map[string]int{
					"p1Wins": m.seriesP1Wins,
					"p2Wins": m.seriesP2Wins,
					"draws":  m.seriesDraws,
				},
			}),
		})

		m.broadcastLocked(Envelope{Type: "game_finished", Payload: mustJSON(map[string]string{"winner": m.winner})})
		m.broadcastStateLocked()
		m.persistLocked()
		return
	}

	// сразу стартуем следующий раунд (как ты хотел)
	m.persistLocked()
	m.startRoundLocked()
}

func (m *Match) attemptLocked(slot Slot) Attempt {
	var me *Player
	var opp *Player
	if slot == P1 {
		me = m.p1
		opp = m.p2
	} else {
		me = m.p2
		opp = m.p1
	}

	// пропуск
	if me.missed || !me.guessSet {
		return Attempt{
			Guess:  nil,
			Bulls:  0,
			Cows:   0,
			Missed: true,
		}
	}

	g := me.guess
	b, c := BullsCows(opp.secret, g)
	return Attempt{
		Guess:  &g,
		Bulls:  b,
		Cows:   c,
		Missed: false,
	}
}

func (m *Match) buildStateLocked(slot Slot) StatePayload {
	you := string(slot)

	connected := 0
	if m.p1.connected {
		connected++
	}
	if m.p2.connected {
		connected++
	}

	return StatePayload{
		MatchID:          m.id,
		You:              you,
		PlayersConnected: connected,
		Phase:            m.phase,
		Round:            m.round,
		DeadlineMs:       toMs(m.deadline),
		SecretsReady: map[string]bool{
			"p1": m.p1.secretSet,
			"p2": m.p2.secretSet,
		},
		GuessesReady: map[string]bool{
			"p1": m.p1.guessSet || m.p1.missed,
			"p2": m.p2.guessSet || m.p2.missed,
		},
		History: m.history,
		Winner:  m.winner,
	}
}

func (m *Match) playerLocked(slot Slot) *Player {
	if slot == P1 {
		return m.p1
	}
	return m.p2
}

func (m *Match) sendLocked(conn *ClientConn, env Envelope) {
	if conn == nil {
		return
	}
	b, _ := json.Marshal(env)
	select {
	case conn.send <- b:
	default:
		// MVP: если клиент не успевает читать, просто дропаем (позже сделаем backpressure)
	}
}

func (m *Match) broadcastLocked(env Envelope) {
	if m.p1.conn != nil {
		m.sendLocked(m.p1.conn, env)
	}
	if m.p2.conn != nil {
		m.sendLocked(m.p2.conn, env)
	}
}

func valid4Digits(s string) bool {
	if len(s) != 4 {
		return false
	}
	for i := 0; i < 4; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func toMs(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

func (m *Match) persistLocked() {
	if m.onPersist == nil {
		return
	}
	m.onPersist(m.snapshotLocked())
}
