package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type PokerClass struct {
	Deck []Card
}

func (p *PokerClass) NewDeck() {
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}

	for _, suit := range suits {
		for _, rank := range ranks {
			p.Deck = append(p.Deck, Card{Rank: rank, Suit: suit, Value: p.CardsIndex(rank, suit)})
		}
	}
}

func (p *PokerClass) Shuffle() {
	n := len(p.Deck) - 1
	for n > 0 {
		r := rand.Float64()
		b := math.Floor(r * float64(n))

		a := p.Deck[n]
		p.Deck[n] = p.Deck[int(b)]
		p.Deck[(int(b))] = a
		n--
	}
}

func (p *PokerClass) Deal(numberOfPlayers int) Table {
	players := [][]Card{}

	for i := 0; i < numberOfPlayers; i++ {
		players = append(players, []Card{p.Deck[i], p.Deck[i+numberOfPlayers]})
	}

	return Table{players, TableDeck{Flop: []Card{p.Deck[numberOfPlayers*2+1], p.Deck[numberOfPlayers*2+2], p.Deck[numberOfPlayers*2+3]}, Turn: p.Deck[numberOfPlayers*2+5], River: p.Deck[numberOfPlayers*2+7]}}
}

func (p *PokerClass) CardsIndex(rank string, suit string) int32 {
	var defaultRankMap map[string]int32 = map[string]int32{"2": 0, "3": 1, "4": 2, "5": 3, "6": 4, "7": 5, "8": 6, "9": 7, "T": 8, "J": 9, "Q": 10, "K": 11, "A": 12}
	var defaultSuitMap map[string]int32 = map[string]int32{"Clubs": 0, "Diamonds": 1, "Hearts": 2, "Spades": 3}
	return 13*defaultSuitMap[suit] + defaultRankMap[rank]
}

func (p *PokerClass) GetAllTableIndexedWins(table Table, lookupTable [32487834]int32) {
	var playersResult []int32

	// fmt.Println("Table in GetAllTableIndexedWins: ", table.Deck, "\n")
	for _, player := range table.Players {
		// fmt.Println("Player hand in GetAllTableIndexedWins: ", player[0], " ", player[1], "\n")
		hand := []int32{
			player[0].Value, player[1].Value,
			table.Deck.Flop[0].Value, table.Deck.Flop[1].Value, table.Deck.Flop[2].Value,
			table.Deck.Turn.Value,
			table.Deck.River.Value,
		}

		fmt.Println("Evaluating hand:", hand)
		idx := int32(53)
		for _, c := range hand {
			fmt.Printf("lookup[%d + %d] = ", idx, c)
			idx = lookupTable[idx+c]
			fmt.Printf("%d\n", idx)
		}
		final := lookupTable[idx]
		fmt.Println("Final rank:", final)

		var index int32 = 53
		for _, cards := range hand {
			index = lookupTable[index+cards]
		}
		rank := lookupTable[index]
		playersResult = append(playersResult, rank)
	}

	fmt.Println(playersResult)
}
