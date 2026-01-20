package game

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/bc-mvp/internal/auth"
	"github.com/gorilla/websocket"
)

type memPersist struct {
	m map[string]MatchSnapshot
}

func (p *memPersist) Save(ctx context.Context, matchID string, snap MatchSnapshot) error {
	if p.m == nil {
		p.m = make(map[string]MatchSnapshot)
	}
	p.m[matchID] = snap
	return nil
}

func (p *memPersist) Load(ctx context.Context, matchID string) (MatchSnapshot, bool, error) {
	snap, ok := p.m[matchID]
	return snap, ok, nil
}

type testVerifier struct{}

func (v testVerifier) Verify(token string) (*auth.Claims, error) {
	if token != "good" {
		return nil, errors.New("bad token")
	}
	return &auth.Claims{UserID: "u1", DisplayName: "Alice"}, nil
}

func TestWS_Endpoint_PathParam(t *testing.T) {
	cfg := Config{RoundDuration: 0}
	persist := &memPersist{}
	matchSvc := NewMatchService(cfg, persist)
	server := NewServer(cfg, matchSvc, testVerifier{})

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mkWSURL := func(path string) string {
		return "ws" + strings.TrimPrefix(ts.URL, "http") + path
	}

	// create one match for happy paths
	const matchID = "abc123"
	if _, err := matchSvc.Create(context.Background(), matchID); err != nil {
		t.Fatalf("create match: %v", err)
	}

	cases := []struct {
		name        string
		urlPath     string
		authHeader  bool
		sendAuthMsg bool
		wantCode    int // 0 => expect success (101)
	}{
		{name: "success_auth_header", urlPath: "/ws/" + matchID, authHeader: true, wantCode: 0},
		{name: "success_auth_message", urlPath: "/ws/" + matchID, sendAuthMsg: true, wantCode: 0},
		{name: "success_ignores_query", urlPath: "/ws/" + matchID + "?matchId=wrong", sendAuthMsg: true, wantCode: 0},
		{name: "bad_missing", urlPath: "/ws/", wantCode: http.StatusBadRequest},
		{name: "bad_extra_segment", urlPath: "/ws/" + matchID + "/x", wantCode: http.StatusBadRequest},
		{name: "not_found", urlPath: "/ws/unknown", wantCode: http.StatusNotFound},
		{name: "unauthorized_header", urlPath: "/ws/" + matchID, authHeader: true, wantCode: http.StatusUnauthorized},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dialer := websocket.Dialer{}
			hdr := http.Header{}
			if tc.authHeader {
				// for unauthorized_header case we pass a bad token
				tok := "good"
				if tc.name == "unauthorized_header" {
					tok = "bad"
				}
				hdr.Set("Authorization", "Bearer "+tok)
			}

			ws, resp, err := dialer.Dial(mkWSURL(tc.urlPath), hdr)
			if tc.wantCode != 0 {
				if err == nil {
					_ = ws.Close()
					t.Fatalf("expected dial error, got nil")
				}
				if resp == nil {
					t.Fatalf("expected HTTP response with status %d, got nil resp (err=%v)", tc.wantCode, err)
				}
				if resp.StatusCode != tc.wantCode {
					t.Fatalf("status=%d, want %d (err=%v)", resp.StatusCode, tc.wantCode, err)
				}
				return
			}

			if err != nil {
				code := 0
				if resp != nil {
					code = resp.StatusCode
				}
				t.Fatalf("dial: status=%d err=%v", code, err)
			}
			defer ws.Close()

			if tc.sendAuthMsg {
				_ = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth","payload":{"token":"good"}}`))
			}

			// wait for a state message
			_ = ws.SetReadDeadline(time.Now().Add(2 * time.Second))
			for {
				_, data, rerr := ws.ReadMessage()
				if rerr != nil {
					t.Fatalf("read: %v", rerr)
				}
				var env Envelope
				if jerr := json.Unmarshal(data, &env); jerr != nil {
					continue
				}
				if env.Type == "state" {
					return
				}
			}
		})
	}
}
