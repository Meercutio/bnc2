package game

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newTestConn() *ClientConn {
	return &ClientConn{
		ws:   nil,
		send: make(chan []byte, 256),
	}
}

func readEnvelopesNonBlocking(c *ClientConn) []Envelope {
	var envs []Envelope
	for {
		select {
		case msg := <-c.send:
			var env Envelope
			if json.Unmarshal(msg, &env) == nil {
				envs = append(envs, env)
			}
		default:
			return envs
		}
	}
}

func findLastState(envs []Envelope) (StatePayload, bool) {
	for i := len(envs) - 1; i >= 0; i-- {
		if envs[i].Type != "state" {
			continue
		}
		var st StatePayload
		if json.Unmarshal(envs[i].Payload, &st) == nil {
			return st, true
		}
	}
	return StatePayload{}, false
}

func TestMatch_Scenarios(t *testing.T) {
	type scenario struct {
		name string
		run  func(t *testing.T)
	}

	cases := []scenario{
		{
			name: "start round after both secrets (roundActive=true, round=1, phase=playing)",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				m.Attach("u1", "Alice", newTestConn())
				m.Attach("u2", "Bob", newTestConn())

				if err := m.SetSecret(P1, "0011"); err != nil {
					t.Fatalf("SetSecret P1: %v", err)
				}
				if err := m.SetSecret(P2, "0101"); err != nil {
					t.Fatalf("SetSecret P2: %v", err)
				}

				m.mu.Lock()
				defer m.mu.Unlock()

				if m.phase != "playing" {
					t.Fatalf("phase=%s want playing", m.phase)
				}
				if m.round != 1 {
					t.Fatalf("round=%d want 1", m.round)
				}
				if !m.roundActive {
					t.Fatalf("roundActive=false want true")
				}
			},
		},
		{
			name: "both submit -> history appended and draw if both guessed",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				m.Attach("u1", "Alice", newTestConn())
				m.Attach("u2", "Bob", newTestConn())

				_ = m.SetSecret(P1, "0011")
				_ = m.SetSecret(P2, "0101")

				if err := m.SubmitGuess(P1, "0101"); err != nil {
					t.Fatalf("SubmitGuess P1: %v", err)
				}
				if err := m.SubmitGuess(P2, "0011"); err != nil {
					t.Fatalf("SubmitGuess P2: %v", err)
				}

				m.mu.Lock()
				defer m.mu.Unlock()

				if len(m.history) != 1 {
					t.Fatalf("history len=%d want 1", len(m.history))
				}
				if m.phase != "finished" {
					t.Fatalf("phase=%s want finished", m.phase)
				}
				if m.winner != "draw" {
					t.Fatalf("winner=%s want draw", m.winner)
				}

				got := m.history[0]
				if got.P1.Guess == nil || *got.P1.Guess != "0101" || got.P1.Bulls != 4 {
					t.Fatalf("P1 attempt unexpected: %+v", got.P1)
				}
				if got.P2.Guess == nil || *got.P2.Guess != "0011" || got.P2.Bulls != 4 {
					t.Fatalf("P2 attempt unexpected: %+v", got.P2)
				}
			},
		},
		{
			name: "winner p1 when p1 guesses opponent secret and p2 doesn't",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				m.Attach("u1", "Alice", newTestConn())
				m.Attach("u2", "Bob", newTestConn())

				_ = m.SetSecret(P1, "9999")
				_ = m.SetSecret(P2, "1234")

				_ = m.SubmitGuess(P1, "1234") // win
				_ = m.SubmitGuess(P2, "0000") // no win

				m.mu.Lock()
				defer m.mu.Unlock()

				if m.phase != "finished" {
					t.Fatalf("phase=%s want finished", m.phase)
				}
				if m.winner != "p1" {
					t.Fatalf("winner=%s want p1", m.winner)
				}
			},
		},
		{
			name: "timeout marks missed and can still finish game (p1 wins, p2 missed)",
			run: func(t *testing.T) {
				m := NewMatch("m1", 50*time.Millisecond)
				m.Attach("u1", "Alice", newTestConn())
				m.Attach("u2", "Bob", newTestConn())

				_ = m.SetSecret(P1, "1111")
				_ = m.SetSecret(P2, "2222")

				_ = m.SubmitGuess(P1, "2222") // p1 guessed opponent

				time.Sleep(90 * time.Millisecond) // > roundDur to avoid flake

				m.mu.Lock()
				defer m.mu.Unlock()

				if len(m.history) < 1 {
					t.Fatalf("expected at least 1 history item")
				}
				last := m.history[len(m.history)-1]
				if !last.P2.Missed || last.P2.Guess != nil {
					t.Fatalf("expected P2 missed with nil guess, got %+v", last.P2)
				}
				if m.phase != "finished" || m.winner != "p1" {
					t.Fatalf("expected finished winner=p1, got phase=%s winner=%s", m.phase, m.winner)
				}
			},
		},
		{
			name: "cannot attach third player (match_full)",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				_, code, _ := m.Attach("u1", "Alice", newTestConn())
				if code != "" {
					t.Fatalf("unexpected code for u1: %s", code)
				}
				_, code, _ = m.Attach("u2", "Bob", newTestConn())
				if code != "" {
					t.Fatalf("unexpected code for u2: %s", code)
				}
				_, code, _ = m.Attach("u3", "Charlie", newTestConn())
				if code != "match_full" {
					t.Fatalf("expected match_full, got %q", code)
				}
			},
		},
		{
			name: "state contains correct you field per connection",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				c1 := newTestConn()
				c2 := newTestConn()

				slot1, code, _ := m.Attach("u1", "Alice", c1)
				if code != "" || slot1 != P1 {
					t.Fatalf("attach u1 => slot=%s code=%s", slot1, code)
				}
				slot2, code, _ := m.Attach("u2", "Bob", c2)
				if code != "" || slot2 != P2 {
					t.Fatalf("attach u2 => slot=%s code=%s", slot2, code)
				}

				// просим отправить персональный state
				m.SendStateTo(P1)
				m.SendStateTo(P2)

				envs1 := readEnvelopesNonBlocking(c1)
				envs2 := readEnvelopesNonBlocking(c2)

				st1, ok := findLastState(envs1)
				if !ok {
					t.Fatalf("no state for p1")
				}
				st2, ok := findLastState(envs2)
				if !ok {
					t.Fatalf("no state for p2")
				}

				if st1.You != "p1" {
					t.Fatalf("st1.You=%s want p1", st1.You)
				}
				if st2.You != "p2" {
					t.Fatalf("st2.You=%s want p2", st2.You)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}

func TestMatch_Scenarios2(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "round starts after both secrets",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)

				_, code, _ := m.Attach("u1", "Alice", newTestConn())
				require.Empty(t, code)

				_, code, _ = m.Attach("u2", "Bob", newTestConn())
				require.Empty(t, code)

				require.NoError(t, m.SetSecret(P1, "0011"))
				require.NoError(t, m.SetSecret(P2, "0101"))

				m.mu.Lock()
				defer m.mu.Unlock()

				assert.Equal(t, "playing", m.phase)
				assert.Equal(t, 1, m.round)
				assert.True(t, m.roundActive)
			},
		},
		{
			name: "p1 wins when guessing opponent secret",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				m.Attach("u1", "Alice", newTestConn())
				m.Attach("u2", "Bob", newTestConn())

				require.NoError(t, m.SetSecret(P1, "9999"))
				require.NoError(t, m.SetSecret(P2, "1234"))

				require.NoError(t, m.SubmitGuess(P1, "1234"))
				require.NoError(t, m.SubmitGuess(P2, "0000"))

				m.mu.Lock()
				defer m.mu.Unlock()

				assert.Equal(t, "finished", m.phase)
				assert.Equal(t, "p1", m.winner)
			},
		},
		{
			name: "timeout marks missed",
			run: func(t *testing.T) {
				m := NewMatch("m1", 40*time.Millisecond)
				m.Attach("u1", "Alice", newTestConn())
				m.Attach("u2", "Bob", newTestConn())

				require.NoError(t, m.SetSecret(P1, "1111"))
				require.NoError(t, m.SetSecret(P2, "2222"))

				require.NoError(t, m.SubmitGuess(P1, "2222"))
				time.Sleep(70 * time.Millisecond)

				m.mu.Lock()
				defer m.mu.Unlock()

				require.NotEmpty(t, m.history)
				last := m.history[len(m.history)-1]

				assert.True(t, last.P2.Missed)
				assert.Equal(t, "finished", m.phase)
				assert.Equal(t, "p1", m.winner)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}

func TestMatch_State_PlayerNames_And_RevealedSecrets(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "playerNames are included in state after attach",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				c1 := newTestConn()
				c2 := newTestConn()
				m.Attach("u1", "Alice", c1)
				m.Attach("u2", "Bob", c2)

				m.SendStateTo(P1)
				st, ok := findLastState(readEnvelopesNonBlocking(c1))
				require.True(t, ok)
				require.Equal(t, map[string]string{"p1": "Alice", "p2": "Bob"}, st.PlayerNames)
				require.Nil(t, st.RevealedSecrets)
			},
		},
		{
			name: "revealedSecrets are present only after finished",
			run: func(t *testing.T) {
				m := NewMatch("m1", 0)
				c1 := newTestConn()
				c2 := newTestConn()
				m.Attach("u1", "Alice", c1)
				m.Attach("u2", "Bob", c2)

				require.NoError(t, m.SetSecret(P1, "1111"))
				require.NoError(t, m.SetSecret(P2, "2222"))
				require.NoError(t, m.SubmitGuess(P1, "2222")) // p1 wins
				require.NoError(t, m.SubmitGuess(P2, "0000"))

				m.SendStateTo(P1)
				st, ok := findLastState(readEnvelopesNonBlocking(c1))
				require.True(t, ok)
				require.Equal(t, "finished", st.Phase)
				require.Equal(t, map[string]string{"p1": "Alice", "p2": "Bob"}, st.PlayerNames)
				require.Equal(t, map[string]string{"p1": "1111", "p2": "2222"}, st.RevealedSecrets)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}
