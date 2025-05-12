package player

import (
	"clyde1811/dmp/cardset"
	"clyde1811/dmp/crypto"
	"fmt"
	"math/big"
	"math/rand"

	eciesgo "github.com/ecies/go/v2"
)

type Player struct {
	Id         int
	Hand       []cardset.CardPoint
	Points     int
	Active     bool
	PublicKey  crypto.CurvePoint
	PrivateKey *big.Int
}

type GameInfo interface {
	CanPlayCard(card cardset.CardPoint) bool
	IsEffectCard(card cardset.CardPoint) bool
	CardValue(card cardset.CardPoint) int
}

func GenerateAsymmetricKey() (*eciesgo.PublicKey, *eciesgo.PrivateKey, error) {
	privKey, err := eciesgo.GenerateKey()

	if err != nil {
		return nil, nil, err
	}

	return privKey.PublicKey, privKey, nil

}

func GeneratePlayer(context crypto.CryptoContext, privKey *big.Int, id int) Player {
	publicKey := context.GeneratePublicKey(privKey)

	return Player{id, []cardset.CardPoint{}, 0, true, publicKey, privKey}
}

func (p *Player) GetPublicKey() crypto.CurvePoint {
	return p.PublicKey
}

func (p *Player) GetPrivateKey() *big.Int {
	return p.PrivateKey
}

func (p *Player) EncryptCard(context crypto.CryptoContext, cardset cardset.CardSet, r *big.Int) {
	cardset.RemaskAllCard(context, r)
}

func (p *Player) DecryptCard(context crypto.CryptoContext, card cardset.Card) cardset.Card {
	Cb := context.DecryptCard(card.Ca, card.Cb, p.PrivateKey)
	card.Cb = Cb

	return card
}

func (p *Player) EstablishCard(card cardset.Card, cardset cardset.CardSet) cardset.CardPoint {
	return cardset.FindCardByPoint(card.Cb)
}

func (p *Player) ChooseCard(g GameInfo) (cardset.CardPoint, int) {
	playable := []int{}
	for i, card := range p.Hand {
		if g.CanPlayCard(card) { // Valid card
			playable = append(playable, i)
		}
	}

	// Check if there is a card can be give out
	if len(playable) == 0 {
		return cardset.CardPoint{}, -1
	}

	// Strategy:
	// 1. Prefer valid non-effect card with the largest value
	// 2. Otherwise, random effect card

	idx := playable[rand.Intn(len(playable))]
	maxVal := -1
	maxIdx := -1
	for _, handIdx := range playable {
		if g.IsEffectCard(p.Hand[handIdx]) != true {
			val := g.CardValue(p.Hand[handIdx])
			if maxIdx == -1 || maxVal < val {
				maxIdx = handIdx
				maxVal = val
			}
		}
	}
	if maxIdx != -1 {
		idx = maxIdx
	}
	card := p.Hand[idx]

	// Remove card from Hand
	p.Hand = append(p.Hand[:idx], p.Hand[idx+1:]...)

	return card, idx
}

func (p *Player) ShowHand() {
	fmt.Printf("[Player %d] Hand:\n", p.Id)
	for i := 0; i < len(p.Hand); i++ {
		c := p.Hand[i]
		fmt.Printf("%d, %s\n", i, c.ToCardString())
	}

	fmt.Println("---")
}
