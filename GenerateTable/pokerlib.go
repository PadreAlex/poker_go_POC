package main

func min(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

var HandCategoryNames = []string{
	"BAD!!",
	"High Card",
	"Pair",
	"Two Pair",
	"Three of a Kind",
	"Straight",
	"Flush",
	"Full House",
	"Four of a Kind",
	"Straight Flush",
}

var IDs [612978]int64
var HR [32487834]int

var numIDs int = 1
var numcards int = 0
var maxHR int = 0
var maxID int64 = 0

func swapIfLess(workcards *[8]int, i, j int) {
	if workcards[i] < workcards[j] {
		workcards[i], workcards[j] = workcards[j], workcards[i]
	}
}

func MakeID(IDin int64, newcard int) int64 {
	var ID int64 = 0
	var suitcount [4 + 1]int
	var rankcount [13 + 1]int
	var workcards [8]int // intentially keeping one as a 0 end
	var cardnum int
	var getout int = 0

	for cardnum := 0; cardnum < 6; cardnum++ {
		workcards[cardnum+1] = int((IDin >> (8 * cardnum)) & 0xff)
	}

	newcard--

	workcards[0] = (((newcard >> 2) + 1) << 4) + (newcard & 3) + 1

	for _, card := range workcards {
		if card == 0 {
			break
		}

		suit := card & 0xf
		rank := (card >> 4) & 0xf

		suitcount[suit]++
		rankcount[rank]++

		if numcards > 0 && workcards[0] == card {
			getout = 1
		}

		numcards++
	}

	if getout == 1 {
		return 0
	}

	var needsuited int = numcards - 2

	if numcards > 4 {
		for rank := 0; rank < 14; rank++ {
			if rankcount[rank] > 4 {
				return 0
			}
		}
	}

	if needsuited > 1 {
		for cardnum = 0; cardnum < numcards; cardnum++ { // for each card
			if suitcount[workcards[cardnum]&0xf] < needsuited { // check suitcount to the number I need to have suits significant
				workcards[cardnum] &= 0xf0 // if not enough - 0 out the suit - now this suit would be a 0 vs 1-4
			}
		}
	}

	swapIfLess(&workcards, 0, 4)
	swapIfLess(&workcards, 1, 5)
	swapIfLess(&workcards, 2, 6)
	swapIfLess(&workcards, 0, 2)
	swapIfLess(&workcards, 1, 3)
	swapIfLess(&workcards, 4, 6)
	swapIfLess(&workcards, 2, 4)
	swapIfLess(&workcards, 3, 5)
	swapIfLess(&workcards, 0, 1)
	swapIfLess(&workcards, 2, 3)
	swapIfLess(&workcards, 4, 5)
	swapIfLess(&workcards, 1, 4)
	swapIfLess(&workcards, 3, 6)
	swapIfLess(&workcards, 1, 2)
	swapIfLess(&workcards, 3, 4)
	swapIfLess(&workcards, 5, 6)

	ID = int64(workcards[0]) +
		(int64(workcards[1]) << 8) +
		(int64(workcards[2]) << 16) +
		(int64(workcards[3]) << 24) +
		(int64(workcards[4]) << 32) +
		(int64(workcards[5]) << 40) +
		(int64(workcards[6]) << 48)

	return ID
}

func SaveID(ID int64) int {
	if ID == 0 {
		return 0
	}

	if ID >= maxID {
		if ID > maxID {
			IDs[numIDs] = ID
			numIDs++
			maxID = ID
		}
		return numIDs - 1
	}

	low := 0
	high := numIDs - 1

	for high-low > 1 {
		mid := (high + low + 1) / 2
		diff := IDs[mid] - ID
		if diff > 0 {
			high = mid
		} else if diff < 0 {
			low = mid
		} else {
			return mid
		}
	}

	copy(IDs[high+1:], IDs[high:numIDs])
	IDs[high] = ID
	numIDs++

	return high
}

func DoEval(IDin int64) int {
	var handrank int = 0
	var cardnum int
	var workcard int
	var rank int
	var suit int
	var mainsuit int = 20
	var suititerator int = 1
	var holdrank int
	var workcards [8]int
	var holdcards [8]int
	var numevalcards int = 0

	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41}
	if IDin != 0 {

	}
}
