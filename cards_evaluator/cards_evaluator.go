package cards_evaluator

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	poker_helper "poker/helpers"
	poker_structs "poker/structs"
)

type PokerClass struct {
	Deck  []poker_structs.Card
	TABLE [32487834]int32
}

func (p *PokerClass) GetTable() {
	file, err := os.Open("./HandRanks.dat")
	poker_helper.CheckError(err)
	defer file.Close()

	stat, err := file.Stat()
	poker_helper.CheckError(err)

	expectedSize := int64(32487834 * 4)
	if stat.Size() != expectedSize {
		panic(fmt.Sprintf("HandRanks.dat is the wrong size: expected %d bytes", expectedSize))
	}

	err = binary.Read(file, binary.LittleEndian, &p.TABLE)
	poker_helper.CheckError(err)
}

func (p *PokerClass) NewDeck() {
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	counter := 1
	for _, rank := range ranks {
		for _, suit := range suits {
			p.Deck = append(p.Deck, poker_structs.Card{Rank: rank, Suit: suit, Value: int32(counter)})
			counter++
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

func (p *PokerClass) Deal(numberOfPlayers int) poker_structs.Table {
	players := [][]poker_structs.Card{}

	for i := range numberOfPlayers {
		players = append(players, []poker_structs.Card{p.Deck[i], p.Deck[i+numberOfPlayers]})
	}

	return poker_structs.Table{Players: players, Deck: poker_structs.TableDeck{Flop: []poker_structs.Card{p.Deck[numberOfPlayers*2+1], p.Deck[numberOfPlayers*2+2], p.Deck[numberOfPlayers*2+3]}, Turn: p.Deck[numberOfPlayers*2+5], River: p.Deck[numberOfPlayers*2+7]}}
}

func (p *PokerClass) LookupHand(cards []int32) int32 {
	if len(cards) != 7 {
		panic(fmt.Sprintf("7 cards lookup. Wrong cards length: %d", len(cards)))
	}
	eval := p.TABLE[53+cards[0]]
	eval = p.TABLE[eval+cards[1]]
	eval = p.TABLE[eval+cards[2]]
	eval = p.TABLE[eval+cards[3]]
	eval = p.TABLE[eval+cards[4]]
	eval = p.TABLE[eval+cards[5]]
	return p.TABLE[eval+cards[6]]
}

func (p *PokerClass) GetAllTableIndexedWins(table poker_structs.Table) []poker_structs.TableResults {
	var playersResult []poker_structs.TableResults

	for i, player := range table.Players {
		hand := []int32{
			player[0].Value, player[1].Value,
			table.Deck.Flop[0].Value, table.Deck.Flop[1].Value, table.Deck.Flop[2].Value,
			table.Deck.Turn.Value,
			table.Deck.River.Value,
		}

		for _, v := range hand {
			if v < 0 || v > 52 {
				panic(fmt.Sprintf("Invalid card value: %d", v))
			}
		}

		eval := p.LookupHand(hand)
		playersResult = append(playersResult, poker_structs.TableResults{Salt: eval & 0x00000FFF, Type: poker_structs.HandRanks[eval>>12], BaseValue: eval, PlayerNumber: i})
	}
	return playersResult
}
