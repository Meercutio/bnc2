package game

func BullsCows(secret, guess string) (bulls, cows int) {
	// bulls
	usedS := [4]bool{}
	usedG := [4]bool{}

	for i := 0; i < 4; i++ {
		if secret[i] == guess[i] {
			bulls++
			usedS[i] = true
			usedG[i] = true
		}
	}

	// counts for remaining
	var cntS [10]int
	var cntG [10]int

	for i := 0; i < 4; i++ {
		if !usedS[i] {
			cntS[int(secret[i]-'0')]++
		}
		if !usedG[i] {
			cntG[int(guess[i]-'0')]++
		}
	}

	for d := 0; d < 10; d++ {
		if cntS[d] < cntG[d] {
			cows += cntS[d]
		} else {
			cows += cntG[d]
		}
	}

	return bulls, cows
}
