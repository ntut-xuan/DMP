package cardset

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Card struct {
	Suite        string
	Number       string
	RandomNumber int64
}

type CardSet struct {
	Rand  *rand.Rand
	Card  [52]Card
	Index int
}

func CreateCardSet(seed int64) *CardSet {
	var card [52]Card
	var suite = [4]string{"H", "S", "C", "D"}
	var code = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var r = rand.New(rand.NewSource(seed))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 13; j++ {
			card[i*13+j-1] = Card{
				Suite:        suite[i],
				Number:       code[j-1],
				RandomNumber: r.Int63(),
			}
		}
	}

	return &CardSet{Rand: r, Card: card, Index: 0}
}

func ParseCardByString(cardString string) (Card, error) {
	cardSegments := strings.Split(cardString, "_")

	if len(cardSegments) != 3 {
		return Card{}, errors.New("ParseCardByString: Failed to parse card. Segment is not equal to 3.")
	}

	randomNumber, err := strconv.ParseInt(cardSegments[2][2:], 16, 64)

	if err != nil {
		return Card{}, err
	}

	return Card{Suite: cardSegments[0], Number: cardSegments[1], RandomNumber: randomNumber}, nil
}

func (cs *CardSet) ShuffleCardSet() {
	cs.Rand.Shuffle(len(cs.Card), func(i, j int) {
		cs.Card[i], cs.Card[j] = cs.Card[j], cs.Card[i]
	})
}

func (cs *CardSet) Draw() Card {
	var index = cs.Index
	var card = cs.Card[index]
	cs.Index += 1
	return card
}

func (c *Card) ToCardString() string {
	return c.Suite + "_" + c.Number + "_" + fmt.Sprintf("0x%016x", c.RandomNumber)
}
