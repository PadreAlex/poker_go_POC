package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

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
	var handrank int
	var holdrank int
	var workcards [8]int
	var holdcards [8]int
	numevalcards := 0
	mainsuit := 20
	suititerator := 1
	primes := [...]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41}

	// Extract cards from the 64-bit ID
	if IDin != 0 {
		for cardnum := 0; cardnum < 7; cardnum++ {
			card := int((IDin >> (8 * cardnum)) & 0xff)
			if card == 0 {
				break
			}
			holdcards[cardnum] = card
			numevalcards++
			if suit := card & 0xf; suit != 0 {
				mainsuit = suit
			}
		}

		// Convert to Cactus Kev's format
		for cardnum := 0; cardnum < numevalcards; cardnum++ {
			workcard := holdcards[cardnum]
			rank := (workcard >> 4) - 1
			suit := workcard & 0xf

			if suit == 0 {
				suit = suititerator
				suititerator++
				if suititerator == 5 {
					suititerator = 1
				}
				if suit == mainsuit {
					suit = suititerator
					suititerator++
					if suititerator == 5 {
						suititerator = 1
					}
				}
			}

			workcards[cardnum] = primes[rank] | (rank << 8) | (1 << (suit + 11)) | (1 << (16 + rank))
		}

		// Evaluate hand
		switch numevalcards {
		case 5:
			holdrank = eval5(workcards[0], workcards[1], workcards[2], workcards[3], workcards[4])
		case 6:
			holdrank = eval5(workcards[0], workcards[1], workcards[2], workcards[3], workcards[4])
			holdrank = min(holdrank, eval5(workcards[0], workcards[1], workcards[2], workcards[3], workcards[5]))
			holdrank = min(holdrank, eval5(workcards[0], workcards[1], workcards[2], workcards[4], workcards[5]))
			holdrank = min(holdrank, eval5(workcards[0], workcards[1], workcards[3], workcards[4], workcards[5]))
			holdrank = min(holdrank, eval5(workcards[0], workcards[2], workcards[3], workcards[4], workcards[5]))
			holdrank = min(holdrank, eval5(workcards[1], workcards[2], workcards[3], workcards[4], workcards[5]))
		case 7:
			holdrank = eval7(workcards[:7])
		default:
			fmt.Printf("Problem with numcards = %d\n", numevalcards)
			return 0
		}

		handrank = 7463 - holdrank

		switch {
		case handrank < 1278:
			handrank = handrank + 4096*1
		case handrank < 4138:
			handrank = handrank - 1277 + 4096*2
		case handrank < 4996:
			handrank = handrank - 4137 + 4096*3
		case handrank < 5854:
			handrank = handrank - 4995 + 4096*4
		case handrank < 5864:
			handrank = handrank - 5853 + 4096*5
		case handrank < 7141:
			handrank = handrank - 5863 + 4096*6
		case handrank < 7297:
			handrank = handrank - 7140 + 4096*7
		case handrank < 7453:
			handrank = handrank - 7296 + 4096*8
		default:
			handrank = handrank - 7452 + 4096*9
		}
	}

	return handrank
}

func main() {
	var IDslot int
	var ID int64
	var card, count int
	var handTypeSum [10]int

	start := time.Now()

	fmt.Println("\nGetting Card IDs!")

	for IDnum := 0; IDs[IDnum] != 0 || IDnum == 0; IDnum++ {
		for card = 1; card < 53; card++ {
			ID = MakeID(IDs[IDnum], card)
			if numcards < 7 {
				SaveID(ID)
			}
		}
		fmt.Printf("\rID - %d", IDnum)
	}

	fmt.Println("\nSetting HandRanks!")

	for IDnum := 0; IDs[IDnum] != 0 || IDnum == 0; IDnum++ {
		for card = 1; card < 53; card++ {
			ID = MakeID(IDs[IDnum], card)

			if numcards < 7 {
				IDslot = SaveID(ID)*53 + 53
			} else {
				IDslot = DoEval(ID)
			}

			maxHR = IDnum*53 + card + 53
			HR[maxHR] = IDslot
		}

		if numcards == 6 || numcards == 7 {
			HR[IDnum*53+53] = DoEval(IDs[IDnum])
		}

		fmt.Printf("\rID - %d", IDnum)
	}

	fmt.Printf("\nNumber IDs = %d\nmaxHR = %d\n", numIDs, maxHR)
	fmt.Printf("Training seconds = %.2f\n", time.Since(start).Seconds())

	start = time.Now()
	var u0, u1, u2, u3, u4, u5 int

	for c0 := 1; c0 < 53; c0++ {
		u0 = HR[53+c0]
		for c1 := c0 + 1; c1 < 53; c1++ {
			u1 = HR[u0+c1]
			for c2 := c1 + 1; c2 < 53; c2++ {
				u2 = HR[u1+c2]
				for c3 := c2 + 1; c3 < 53; c3++ {
					u3 = HR[u2+c3]
					for c4 := c3 + 1; c4 < 53; c4++ {
						u4 = HR[u3+c4]
						for c5 := c4 + 1; c5 < 53; c5++ {
							u5 = HR[u4+c5]
							for c6 := c5 + 1; c6 < 53; c6++ {
								handTypeSum[HR[u5+c6]>>12]++
								count++
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("\nValidation seconds = %.4f\n", time.Since(start).Seconds())

	for i := 0; i <= 9; i++ {
		fmt.Printf("\n%16s = %d", HandCategoryNames[i], handTypeSum[i])
	}

	fmt.Printf("\nTotal Hands = %d\n", count)

	fout, err := os.Create("HandRanks.dat")
	if err != nil {
		fmt.Println("Problem creating the Output File!")
		os.Exit(1)
	}
	defer fout.Close()

	err = binary.Write(fout, binary.LittleEndian, HR[:])
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
}
