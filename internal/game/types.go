package game

import "encoding/json"

// Envelope WS envelope: {"type":"...","payload":{...}}
type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// SetSecretPayload входящие
type SetSecretPayload struct {
	Secret string `json:"secret"`
}

type SubmitGuessPayload struct {
	Guess string `json:"guess"`
}

// RoundStartedPayload исходящие
type RoundStartedPayload struct {
	Round      int   `json:"round"`
	DeadlineMs int64 `json:"deadlineMs"`
}

type Attempt struct {
	Guess  *string `json:"guess"` // null если пропуск
	Bulls  int     `json:"bulls"`
	Cows   int     `json:"cows"`
	Missed bool    `json:"missed"`
}

type RoundHistoryItem struct {
	Round int     `json:"round"`
	P1    Attempt `json:"p1"`
	P2    Attempt `json:"p2"`
}

type StatePayload struct {
	MatchID          string             `json:"matchId"`
	You              string             `json:"you"` // "p1" | "p2"
	PlayerNames      map[string]string  `json:"playerNames"`
	PlayersConnected int                `json:"playersConnected"`
	Phase            string             `json:"phase"` // waiting_players|waiting_secrets|playing|finished
	Round            int                `json:"round"`
	DeadlineMs       int64              `json:"deadlineMs"`
	SecretsReady     map[string]bool    `json:"secretsReady"` // p1/p2
	GuessesReady     map[string]bool    `json:"guessesReady"` // p1/p2 (текущий раунд)
	History          []RoundHistoryItem `json:"history"`
	Winner           string             `json:"winner"`                    // p1|p2|draw|"" (если не закончено)
	RevealedSecrets  map[string]string  `json:"revealedSecrets,omitempty"` // показываем только после finished
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
