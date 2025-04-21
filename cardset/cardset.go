package cardset

import (
	"fmt"
	"math/rand"
)

type Card struct {
	suite string
	number string
	random_num int64
}

type CardSet struct {
	rand *rand.Rand
	card [52]Card
	index int
}

func CreateCardSet(seed int64) *CardSet {
	var card [52]Card;
	var suite = [4]string{"H", "S", "C", "D"};
	var code = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var r = rand.New(rand.NewSource(seed))

	for i := 0; i < 4; i++ {
		for j := 1; j <= 13; j++ {
			card[i*13+j-1] = Card{
				suite: suite[i],
				number: code[j-1],
				random_num: r.Int63(),
			}
		}
	}

	return &CardSet{rand: r, card: card, index: 0}
}

func (cs *CardSet) ShuffleCardSet() {
	cs.rand.Shuffle(len(cs.card), func(i, j int) {
		cs.card[i], cs.card[j] = cs.card[j], cs.card[i]
	})
}

func (cs *CardSet) Draw() (Card) {
	var index = cs.index
	var card = cs.card[index]
	cs.index += 1
	return card
}

func (c *Card) ToCardString() (string) {
	return c.suite + c.number + " -> " + fmt.Sprintf("0x%016x", c.random_num)
}