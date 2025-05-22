package main

func main() {
	p := PokerClass{}
	l := LookupTable{}
	l.getTable()
	p.NewDeck()
	p.Shuffle()
	// table := p.Deal(4)
	// fmt.Println(table.Players, "\n")
	// fmt.Println(table.Deck.Flop)
	// fmt.Println(table.Deck.Turn)
	// fmt.Println(table.Deck.River, "\n")
	table := Table{
		Players: [][]Card{
			{{Rank: "2", Suit: "Clubs", Value: 0}, {Rank: "3", Suit: "Clubs", Value: 1}},
			{{Rank: "4", Suit: "Clubs", Value: 2}, {Rank: "5", Suit: "Clubs", Value: 3}},
		},
		Deck: TableDeck{
			Flop: []Card{
				{Rank: "A", Suit: "Spades", Value: 51}, // A♠ = 13×3 + 12
				{Rank: "K", Suit: "Spades", Value: 50}, // K♠
				{Rank: "Q", Suit: "Spades", Value: 49}, // Q♠
			},
			Turn:  Card{Rank: "J", Suit: "Spades", Value: 48}, // J♠
			River: Card{Rank: "T", Suit: "Spades", Value: 47}, // T♠
		},
	}

	p.GetAllTableIndexedWins(table, l.TABLE)
}
