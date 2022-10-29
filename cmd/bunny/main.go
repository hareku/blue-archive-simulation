package main

import (
	"fmt"
	"log"
	"math/rand"
)

var r = rand.New(rand.NewSource(100))

type Strategy struct {
	Description string
	Simulate    func(coins int) *SimulateResult
}

type SimulateResult struct {
	ConsumedCoins int
	OpenedCards   int
	Shuffled      int
	Cards         map[CardType]int
}

type CardType int

const (
	CardN CardType = iota
	CardR
	CardSR
	CardUR
)

// SUM: 1828
// N: 564
// R: 414
// SR: 642
// UR: 208
//
// SR: (1828/4)*642/(642+208)=345
// UR: (1828/4)*208/(642+208)=112
func GenerateRandomCards(cards []CardType) {
	gteSRInd := r.Intn(4)

	for i := 0; i < 4; i++ {
		var v CardType
		for {
			f := r.Float32()
			if f < 0.4 {
				v = CardN
			} else if f < 0.7 {
				v = CardR
			} else if f < 0.935 {
				v = CardSR
			} else {
				v = CardUR
			}

			if gteSRInd == i && (v != CardSR && v != CardUR) {
				continue
			}
			break
		}

		cards[i] = CardType(v)
	}
}

var Strategies = []*Strategy{
	{
		Description: "Open All Cards",
		Simulate: func(coins int) *SimulateResult {
			res := &SimulateResult{
				Cards: make(map[CardType]int),
			}
			cards := make([]CardType, 4)
			for res.ConsumedCoins < coins {
				GenerateRandomCards(cards)
				for _, card := range cards {
					res.Cards[card]++
				}
				res.ConsumedCoins += 200 + 210 + 220 + 230
				res.OpenedCards += 4
				res.Shuffled++
			}
			return res
		},
	},

	{
		Description: "Shuffle every card",
		Simulate: func(coins int) *SimulateResult {
			res := &SimulateResult{
				Cards: make(map[CardType]int),
			}
			cards := make([]CardType, 4)
			for res.ConsumedCoins < coins {
				GenerateRandomCards(cards)
				ind := r.Intn(4)

				res.Cards[cards[ind]]++
				res.ConsumedCoins += 200
				res.OpenedCards++
				res.Shuffled++
			}
			return res
		},
	},

	{
		Description: "Shuffle when opened SR or UR",
		Simulate: func(coins int) *SimulateResult {
			res := &SimulateResult{
				Cards: make(map[CardType]int),
			}
			cards := make([]CardType, 4)
			opened := make([]bool, 4)
			for res.ConsumedCoins < coins {
				for i := 0; i < 4; i++ {
					opened[i] = false
				}

				GenerateRandomCards(cards)

				for i := 0; i < 4; i++ {
					var ind int
					for {
						ind = r.Intn(4)
						if !opened[ind] {
							opened[ind] = true
							break
						}
					}

					openedCard := cards[ind]
					res.Cards[openedCard]++
					res.ConsumedCoins += 200 + i*10
					res.OpenedCards++
					if openedCard == CardSR || openedCard == CardUR {
						break
					}
				}

				res.Shuffled++
			}
			return res
		},
	},
}

const Coins = 1_000_000_000

func main() {
	log.Printf("Num of coins to simulate: %d", Coins)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	for _, s := range Strategies {
		runStrategy(s)
	}
	return nil
}

func runStrategy(s *Strategy) {
	res := s.Simulate(Coins)
	log.Printf("%+v", res)

	sumCoins := 0
	for _, v := range res.Cards {
		sumCoins += v
	}

	fmt.Printf("%s\n", s.Description)
	fmt.Printf("N: %+v\n", float64(res.Cards[CardN])/float64(sumCoins))
	fmt.Printf("R: %+v\n", float64(res.Cards[CardR])/float64(sumCoins))
	fmt.Printf("SR: %+v\n", float64(res.Cards[CardSR])/float64(sumCoins))
	fmt.Printf("UR: %+v\n", float64(res.Cards[CardUR])/float64(sumCoins))
	fmt.Printf("Shuffled: %d\n", res.Shuffled)
}
