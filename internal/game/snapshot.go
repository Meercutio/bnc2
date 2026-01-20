package game

import "time"

// MatchSnapshot — сериализуемое состояние матча для Redis.
type MatchSnapshot struct {
	MatchID string `json:"matchId"`

	Phase string `json:"phase"`
	Round int    `json:"round"`

	// важное: сохраняем ID игроков, иначе после рестарта невозможно корректно reconnect
	P1ID   string `json:"p1Id"`
	P1Name string `json:"p1Name,omitempty"`
	P2ID   string `json:"p2Id"`
	P2Name string `json:"p2Name,omitempty"`

	P1Secret    string `json:"p1Secret"`
	P1SecretSet bool   `json:"p1SecretSet"`
	P2Secret    string `json:"p2Secret"`
	P2SecretSet bool   `json:"p2SecretSet"`

	P1Guess    string `json:"p1Guess"`
	P1GuessSet bool   `json:"p1GuessSet"`
	P2Guess    string `json:"p2Guess"`
	P2GuessSet bool   `json:"p2GuessSet"`

	// rematch flags
	P1Rematch bool `json:"p1Rematch"`
	P2Rematch bool `json:"p2Rematch"`

	// счёт серии в рамках matchId
	SeriesP1Wins int `json:"seriesP1Wins"`
	SeriesP2Wins int `json:"seriesP2Wins"`
	SeriesDraws  int `json:"seriesDraws"`

	DeadlineMs int64 `json:"deadlineMs"` // unix millis, 0 если нет дедлайна

	Winner  string             `json:"winner"`
	History []RoundHistoryItem `json:"history"`
}

func (m *Match) snapshotLocked() MatchSnapshot {
	var deadlineMs int64
	if !m.deadline.IsZero() {
		deadlineMs = m.deadline.UnixMilli()
	}

	return MatchSnapshot{
		MatchID: m.id,
		Phase:   m.phase,
		Round:   m.round,

		P1ID:   m.p1.id,
		P1Name: m.p1.name,
		P2ID:   m.p2.id,
		P2Name: m.p2.name,

		P1Secret:    m.p1.secret,
		P1SecretSet: m.p1.secretSet,
		P2Secret:    m.p2.secret,
		P2SecretSet: m.p2.secretSet,

		P1Guess:    m.p1.guess,
		P1GuessSet: m.p1.guessSet,
		P2Guess:    m.p2.guess,
		P2GuessSet: m.p2.guessSet,

		P1Rematch: m.p1.rematchRequested,
		P2Rematch: m.p2.rematchRequested,

		SeriesP1Wins: m.seriesP1Wins,
		SeriesP2Wins: m.seriesP2Wins,
		SeriesDraws:  m.seriesDraws,

		DeadlineMs: deadlineMs,

		Winner:  m.winner,
		History: append([]RoundHistoryItem(nil), m.history...),
	}
}

func (m *Match) restoreLocked(s MatchSnapshot) {
	m.phase = s.Phase
	m.round = s.Round

	// players
	m.p1.id = s.P1ID
	m.p1.name = s.P1Name
	m.p2.id = s.P2ID
	m.p2.name = s.P2Name

	// после рестарта нет соединений
	m.p1.conn = nil
	m.p2.conn = nil
	m.p1.connected = false
	m.p2.connected = false

	// secrets / guesses
	m.p1.secret = s.P1Secret
	m.p1.secretSet = s.P1SecretSet
	m.p2.secret = s.P2Secret
	m.p2.secretSet = s.P2SecretSet

	m.p1.guess = s.P1Guess
	m.p1.guessSet = s.P1GuessSet
	m.p2.guess = s.P2Guess
	m.p2.guessSet = s.P2GuessSet

	// rematch flags
	m.p1.rematchRequested = s.P1Rematch
	m.p2.rematchRequested = s.P2Rematch

	// series score
	m.seriesP1Wins = s.SeriesP1Wins
	m.seriesP2Wins = s.SeriesP2Wins
	m.seriesDraws = s.SeriesDraws

	if s.DeadlineMs > 0 {
		m.deadline = time.UnixMilli(s.DeadlineMs)
	} else {
		m.deadline = time.Time{}
	}

	m.winner = s.Winner
	m.history = append([]RoundHistoryItem(nil), s.History...)

	// активен раунд только если playing
	m.roundActive = (m.phase == "playing")
}
