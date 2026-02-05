package game

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Matchmaker is a simple in-memory matchmaking queue.
//
// MVP limitations:
//   - Works only inside a single server instance (no cross-replica matching).
//   - One-slot queue (pairs two players as soon as the 2nd arrives).
type Matchmaker struct {
	mu   sync.Mutex
	wait *waiter
}

type waiter struct {
	userID      string
	displayName string
	since       time.Time
	ch          chan string // buffered (1) matchId
}

var (
	ErrAlreadyWaiting = errors.New("already waiting")
)

// Join blocks until the user is matched with someone and a matchId is produced,
// or until ctx is cancelled.
//
// createMatch is called exactly once for each successful pair.
func (m *Matchmaker) Join(ctx context.Context, userID, displayName string, createMatch func() (string, error)) (string, error) {
	if userID == "" {
		return "", errors.New("missing userID")
	}
	if createMatch == nil {
		return "", errors.New("missing createMatch")
	}

	// Try to pair with an existing waiter.
	m.mu.Lock()
	if m.wait != nil && m.wait.userID == userID {
		m.mu.Unlock()
		return "", ErrAlreadyWaiting
	}
	if m.wait != nil {
		w := m.wait
		m.wait = nil
		m.mu.Unlock()

		matchID, err := createMatch()
		if err != nil {
			// Best-effort: put the waiter back.
			m.mu.Lock()
			if m.wait == nil {
				m.wait = w
			}
			m.mu.Unlock()
			return "", err
		}

		// Notify first waiter (buffered).
		select {
		case w.ch <- matchID:
		default:
		}
		return matchID, nil
	}

	// Become the waiter.
	w := &waiter{
		userID:      userID,
		displayName: displayName,
		since:       time.Now(),
		ch:          make(chan string, 1),
	}
	m.wait = w
	m.mu.Unlock()

	select {
	case matchID := <-w.ch:
		return matchID, nil
	case <-ctx.Done():
		// Remove self if still waiting.
		m.mu.Lock()
		if m.wait != nil && m.wait.userID == userID {
			m.wait = nil
		}
		m.mu.Unlock()
		return "", ctx.Err()
	}
}

func (m *Matchmaker) Stats() (waiting bool, waitingFor time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.wait == nil {
		return false, 0
	}
	return true, time.Since(m.wait.since)
}
