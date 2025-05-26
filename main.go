package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintln(file, "Error opening or creating file:", err)
		return
	}
	defer file.Close()
	p := PokerClass{}
	p.getTable()
	p.NewDeck()
	for i := 0; i < 100; i++ {
		fmt.Fprintln(file, "Game: ", i, "Starts\n")
		p.Shuffle()
		table := p.Deal(2)
		// fmt.Fprintln(file, "Players: ", table.Players, "\n")
		for plNum, pl := range table.Players {
			fmt.Fprintln(file, "Player", plNum, "With hand: ", pl[0].Rank, pl[0].Suit, pl[1].Rank, pl[1].Suit, "\n")
		}
		fmt.Fprintln(file, "Flop: ", table.Deck.Flop)
		fmt.Fprintln(file, "Turn: ", table.Deck.Turn)
		fmt.Fprintln(file, "River: ", table.Deck.River, "\n")

		results := p.GetAllTableIndexedWins(table)

		winner := results[0]
		for _, result := range results {
			if winner.BaseValue < result.BaseValue {
				winner = result
			}
		}
		fmt.Fprintln(file, results)
		fmt.Fprintln(file, "Winner is player ", winner.PlayerNumber)
		fmt.Fprintln(file, "With hand: ", winner.Type, "\n")
		fmt.Fprintln(file, "Game: ", i, "Ends\n")
	}
	// fmt.Fprintln(file, l.TABLE[50:100])
}
