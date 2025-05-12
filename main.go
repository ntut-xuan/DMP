package main

import (
	"clyde1811/dmp/cardset"
	"clyde1811/dmp/crypto"
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("Decentralized Mental Poker v1.0")
	fmt.Println("We have four players in default.")

	context := crypto.NewContext()

	var priv1 int64 = 31356000600
	var priv2 int64 = 21475841293
	var priv3 int64 = 78123712313
	var priv4 int64 = 61237131292
	var priv5 int64 = 17391831902

	pub1 := context.GeneratePublicKey(new(big.Int).SetInt64(priv1))
	pub2 := context.GeneratePublicKey(new(big.Int).SetInt64(priv2))
	pub3 := context.GeneratePublicKey(new(big.Int).SetInt64(priv3))
	pub4 := context.GeneratePublicKey(new(big.Int).SetInt64(priv4))
	pub5 := context.GeneratePublicKey(new(big.Int).SetInt64(priv5))

	context.GenerateAggregateKey([]crypto.CurvePoint{pub1, pub2, pub3, pub4, pub5})

	// Mask Card
	cardset := cardset.CreateCardSet(12731, *context)

	for i := range 52 {
		fmt.Printf("M: %s (%s, %s)\n", cardset.CardPoint[i].ToCardShrtnString(), cardset.CardPoint[i].Point.X, cardset.CardPoint[i].Point.Y)
	}

	// Player 1 Shuffle
	cardset.ShuffleCardSet()
	cardset.RemaskAllCard(*context, big.NewInt(66739))

	// Player 2 Shuffle
	cardset.ShuffleCardSet()
	cardset.RemaskAllCard(*context, big.NewInt(71283))

	// Player 3 Shuffle
	cardset.ShuffleCardSet()
	cardset.RemaskAllCard(*context, big.NewInt(88123))

	// Player 4 Shuffle
	cardset.ShuffleCardSet()
	cardset.RemaskAllCard(*context, big.NewInt(12573))

	// Player 5 Shuffle
	cardset.ShuffleCardSet()
	cardset.RemaskAllCard(*context, big.NewInt(1))

	// Trying to decrypt
	Ca := cardset.Card[0].Ca
	Cb := cardset.Card[0].Cb

	CbD1 := context.DecryptCard(Ca, Cb, new(big.Int).SetInt64(priv1))
	CbD2 := context.DecryptCard(Ca, CbD1, new(big.Int).SetInt64(priv2))
	CbD3 := context.DecryptCard(Ca, CbD2, new(big.Int).SetInt64(priv3))
	CbD4 := context.DecryptCard(Ca, CbD3, new(big.Int).SetInt64(priv4))
	CbD5 := context.DecryptCard(Ca, CbD4, new(big.Int).SetInt64(priv5))

	fmt.Printf("\nDec M (%s, %s)\n", CbD5.X.String(), CbD5.Y.String())

	for i := range 52 {
		if cardset.CardPoint[i].Point.X.Cmp(CbD5.X) == 0 && cardset.CardPoint[i].Point.Y.Cmp(CbD5.Y) == 0 {
			fmt.Printf("Decrypt Card: %s\n", cardset.CardPoint[i].ToCardShrtnString())
		}
	}
}
