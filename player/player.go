package player

import (
	"clyde1811/dmp/cardset"
	"fmt"
	"math/rand"

	eciesgo "github.com/ecies/go/v2"
)

type Player struct {
	Id         int
	Hand       []cardset.Card
	Points     int
	Active     bool
	PublicKey  *eciesgo.PublicKey
	PrivateKey *eciesgo.PrivateKey
}

type GameInfo interface {
	CanPlayCard(card cardset.Card) bool
	IsEffectCard(card cardset.Card) bool
	CardValue(card cardset.Card) int
}

func GenerateAsymmetricKey() (*eciesgo.PublicKey, *eciesgo.PrivateKey, error) {
	privKey, err := eciesgo.GenerateKey()

	if err != nil {
		return nil, nil, err
	}

	return privKey.PublicKey, privKey, nil

}

func GeneratePlayer(id int) (Player, error) {
	publicKey, privateKey, err := GenerateAsymmetricKey()

	if err != nil {
		return Player{}, err
	}

	return Player{id, []cardset.Card{}, 0, true, publicKey, privateKey}, err
}

func (p *Player) GetPublicKey() string {
	return p.PublicKey.Hex(true)
}

func (p *Player) GetPrivateKey() string {
	return p.PrivateKey.Hex()
}

func EncryptCard(cardPlainText []byte, publicKey *eciesgo.PublicKey) ([]byte, error) {
	ciphertext, err := eciesgo.Encrypt(publicKey, cardPlainText)

	return ciphertext, err
}

func DecryptCard(cardCipherText []byte, privateKey *eciesgo.PrivateKey) ([]byte, error) {
	plaintext, err := eciesgo.Decrypt(privateKey, cardCipherText)

	if err != nil {
		return []byte{}, err
	}

	return plaintext, err
}

func (p *Player) EstablishCard(card cardset.Card) cardset.Card {
	return card
}

func (p *Player) ChooseCard(g GameInfo) (cardset.Card, int) {
	playable := []int{}
	for i, card := range p.Hand {
		if g.CanPlayCard(card) { // Valid card
			playable = append(playable, i)
		}
	}

	// Check if there is a card can be give out
	if len(playable) == 0 {
		return cardset.Card{}, -1
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
