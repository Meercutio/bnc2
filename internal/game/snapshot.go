package game

import "time"

// MatchSnapshot — сериализуемое состояние матча, которое можно положить в Redis.
type MatchSnapshot struct {
	MatchID string `json:"matchId"`

	Phase string `json:"phase"`
	Round int    `json:"round"`

	P1Secret    string `json:"p1Secret"`
	P1SecretSet bool   `json:"p1SecretSet"`
	P2Secret    string `json:"p2Secret"`
	P2SecretSet bool   `json:"p2SecretSet"`

	P1Guess    string `json:"p1Guess"`
	P1GuessSet bool   `json:"p1GuessSet"`
	P2Guess    string `json:"p2Guess"`
	P2GuessSet bool   `json:"p2GuessSet"`

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

		P1Secret:    m.p1.secret,
		P1SecretSet: m.p1.secretSet,
		P2Secret:    m.p2.secret,
		P2SecretSet: m.p2.secretSet,

		P1Guess:    m.p1.guess,
		P1GuessSet: m.p1.guessSet,
		P2Guess:    m.p2.guess,
		P2GuessSet: m.p2.guessSet,

		DeadlineMs: deadlineMs,

		Winner:  m.winner,
		History: append([]RoundHistoryItem(nil), m.history...),
	}
}

func (m *Match) restoreLocked(s MatchSnapshot) {
	m.phase = s.Phase
	m.round = s.Round

	m.p1.secret = s.P1Secret
	m.p1.secretSet = s.P1SecretSet
	m.p2.secret = s.P2Secret
	m.p2.secretSet = s.P2SecretSet

	m.p1.guess = s.P1Guess
	m.p1.guessSet = s.P1GuessSet
	m.p2.guess = s.P2Guess
	m.p2.guessSet = s.P2GuessSet

	if s.DeadlineMs > 0 {
		m.deadline = time.UnixMilli(s.DeadlineMs)
	} else {
		m.deadline = time.Time{}
	}

	m.winner = s.Winner
	m.history = append([]RoundHistoryItem(nil), s.History...)

	// логически: активен, если матч в playing
	m.roundActive = (m.phase == "playing")
}
