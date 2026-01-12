package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBullsCows_Table(t *testing.T) {
	cases := []struct {
		name   string
		secret string
		guess  string
		bulls  int
		cows   int
	}{
		{
			name:   "all match",
			secret: "0011",
			guess:  "0011",
			bulls:  4,
			cows:   0,
		},
		{
			name:   "no match",
			secret: "0000",
			guess:  "1111",
			bulls:  0,
			cows:   0,
		},
		{
			name:   "example with repeats",
			secret: "0011",
			guess:  "0101",
			bulls:  2,
			cows:   2,
		},
		{
			name:   "multiset repeats",
			secret: "1122",
			guess:  "2211",
			bulls:  0,
			cows:   4,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, c := BullsCows(tc.secret, tc.guess)
			assert.Equal(t, tc.bulls, b, "bulls mismatch")
			assert.Equal(t, tc.cows, c, "cows mismatch")
		})
	}
}
