package game

import "time"

// GameFinishedEvent is emitted once per finished game in a match-series.
//
// Winner is one of: p1|p2|draw.
// ResultID is unique for each finished game and is safe for idempotency.
// ResultID format: "<matchId>:<gameNo>".
type GameFinishedEvent struct {
	ResultID   string    `json:"resultId"`
	MatchID    string    `json:"matchId"`
	GameNo     int       `json:"gameNo"`
	P1ID       string    `json:"p1Id"`
	P2ID       string    `json:"p2Id"`
	Winner     string    `json:"winner"`
	FinishedAt time.Time `json:"finishedAt"`
}

type GameResultSink interface {
	HandleGameFinished(e GameFinishedEvent)
}
