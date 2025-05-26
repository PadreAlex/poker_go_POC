package main

type Card struct {
	Rank  string
	Suit  string
	Value int32
}

type Table struct {
	Players [][]Card
	Deck    TableDeck
}

type TableDeck struct {
	Flop  []Card
	Turn  Card
	River Card
}

var HandRanks []string = []string{"BAD!!", "High Card", "Pair", "Two Pair", "Three of a Kind", "Straight", "Flush", "Full House", "Four of a Kind", "Straight Flush"}

type TableResults struct {
	Salt         int32
	Type         string
	BaseValue    int32
	PlayerNumber int
}
