package game

import "testing"

func TestMatchIDFromWSPath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
		want string
		ok   bool
	}{
		{name: "valid", path: "/ws/abc123", want: "abc123", ok: true},
		{name: "valid_longer", path: "/ws/abc123xyz0", want: "abc123xyz0", ok: true},
		{name: "missing", path: "/ws/", want: "", ok: false},
		{name: "missing_no_trailing_slash", path: "/ws", want: "", ok: false},
		{name: "wrong_prefix", path: "/wss/abc", want: "", ok: false},
		{name: "extra_segment", path: "/ws/abc/def", want: "", ok: false},
		{name: "invalid_chars_upper", path: "/ws/Abc", want: "", ok: false},
		{name: "invalid_chars_dash", path: "/ws/abc-def", want: "", ok: false},
		{name: "too_long", path: "/ws/" + makeString('a', 65), want: "", ok: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, ok := matchIDFromWSPath(tc.path)
			if ok != tc.ok {
				t.Fatalf("ok=%v, want %v (got=%q)", ok, tc.ok, got)
			}
			if got != tc.want {
				t.Fatalf("got=%q, want %q", got, tc.want)
			}
		})
	}
}

func makeString(ch byte, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
