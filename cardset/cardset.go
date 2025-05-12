package cardset

import (
	"clyde1811/dmp/crypto"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

type CardPoint struct {
	Suite  string
	Number string
	Point  crypto.CurvePoint
}

type Card struct {
	Ca crypto.CurvePoint
	Cb crypto.CurvePoint
}

type CardSet struct {
	Rand      *rand.Rand
	CardPoint [52]CardPoint
	Card      [52]Card
	Index     int
}

var ErrCardsetEmpty = errors.New("cardset is empty")

func CreateCardSet(seed int64, context crypto.CryptoContext) *CardSet {
	var cardPoint [52]CardPoint
	var card [52]Card
	var suite = [4]string{"H", "S", "C", "D"}
	var code = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var r = rand.New(rand.NewSource(seed))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 13; j++ {
			point := context.GenerateRandomPoint()

			cardPoint[i*13+j-1] = CardPoint{
				Suite:  suite[i],
				Number: code[j-1],
				Point:  point,
			}

			randNumber := r.Int63()
			Ca, Cb := context.MaskCard(point, new(big.Int).SetInt64(randNumber))

			card[i*13+j-1] = Card{
				Ca: Ca,
				Cb: Cb,
			}
		}
	}

	return &CardSet{Rand: r, CardPoint: cardPoint, Card: card, Index: 0}
}

func (cs *CardSet) RemaskAllCard(context crypto.CryptoContext, r *big.Int) {
	for i := range 52 {
		NCa, NCb := context.RemaskCard(cs.Card[i].Ca, cs.Card[i].Cb, r)

		cs.Card[i] = Card{
			Ca: NCa,
			Cb: NCb,
		}
	}
}

func (cs *CardSet) ShuffleCardSet() {
	cs.Rand.Shuffle(len(cs.Card), func(i, j int) {
		cs.Card[i], cs.Card[j] = cs.Card[j], cs.Card[i]
	})
}

func (cs *CardSet) Draw() (Card, error) {
	var index = cs.Index
	if index >= len(cs.Card) {
		return Card{}, ErrCardsetEmpty
	}

	var card = cs.Card[index]
	cs.Index += 1
	return card, nil
}

func (c *CardPoint) ToCardString() string {
	return c.Suite + "_" + c.Number + "_" + fmt.Sprintf("0x%016x", c.Point.X) + "_" + fmt.Sprintf("0x%016x", c.Point.Y)
}

func (c *CardPoint) ToCardShrtnString() string {
	return c.Suite + "_" + c.Number
}
