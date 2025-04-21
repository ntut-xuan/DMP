package player

import (
	"clyde1811/dmp/cardset"

	eciesgo "github.com/ecies/go/v2"
)

type Player struct {
	PublicKey  *eciesgo.PublicKey
	PrivateKey *eciesgo.PrivateKey
}

func GenerateAsymmetricKey() (*eciesgo.PublicKey, *eciesgo.PrivateKey, error) {
	privKey, err := eciesgo.GenerateKey()

	if err != nil {
		return nil, nil, err
	}

	return privKey.PublicKey, privKey, nil

}

func GeneratePlayer() (Player, error) {
	publicKey, privateKey, err := GenerateAsymmetricKey()

	if err != nil {
		return Player{}, err
	}

	return Player{publicKey, privateKey}, err
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
