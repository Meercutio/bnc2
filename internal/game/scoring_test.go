package game

import "testing"

func TestBullsCows_AllMatch(t *testing.T) {
	b, c := BullsCows("0011", "0011")
	if b != 4 || c != 0 {
		t.Fatalf("expected 4 bulls,0 cows got %d bulls,%d cows", b, c)
	}
}

func TestBullsCows_NoMatch(t *testing.T) {
	b, c := BullsCows("0000", "1111")
	if b != 0 || c != 0 {
		t.Fatalf("expected 0,0 got %d,%d", b, c)
	}
}

func TestBullsCows_WithRepeats_Example(t *testing.T) {
	// твой кейс: secret 0011, guess 0101 -> 2 быка, 2 коровы
	b, c := BullsCows("0011", "0101")
	if b != 2 || c != 2 {
		t.Fatalf("expected 2,2 got %d,%d", b, c)
	}
}

func TestBullsCows_RepeatsCountedAsMultiset(t *testing.T) {
	// secret 1122, guess 2211 -> 0 bulls, 4 cows
	b, c := BullsCows("1122", "2211")
	if b != 0 || c != 4 {
		t.Fatalf("expected 0,4 got %d,%d", b, c)
	}
}

func TestValid4Digits(t *testing.T) {
	cases := []struct {
		s  string
		ok bool
	}{
		{"0000", true},
		{"0123", true},
		{"9999", true},
		{"123", false},
		{"12345", false},
		{"12a4", false},
		{"-123", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := valid4Digits(tc.s); got != tc.ok {
			t.Fatalf("valid4Digits(%q)=%v want %v", tc.s, got, tc.ok)
		}
	}
}
