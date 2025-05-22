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
