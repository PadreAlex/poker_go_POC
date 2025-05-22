package main

type HandRank int

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	Trips
	Straight
	Flush
	FullHouse
	Quads
	StraightFlush
	RoyalFlush
)

type Hand struct {
	Rank  HandRank
	Cards []Card
}

func (p *PokerClass) evaluateHandRank() {

}
