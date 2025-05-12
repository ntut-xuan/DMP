package main

import (
	"clyde1811/dmp/crypto"
	"clyde1811/dmp/game"
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("Decentralized Mental Poker v1.0")
	fmt.Println("We have four players in default.")

	context := crypto.NewContext()

	var priv1 *big.Int = new(big.Int).SetInt64(31356000600)
	var priv2 *big.Int = new(big.Int).SetInt64(21475841293)
	var priv3 *big.Int = new(big.Int).SetInt64(78123712313)
	var priv4 *big.Int = new(big.Int).SetInt64(61237131292)

	pub1 := context.GeneratePublicKey(priv1)
	pub2 := context.GeneratePublicKey(priv2)
	pub3 := context.GeneratePublicKey(priv3)
	pub4 := context.GeneratePublicKey(priv4)

	context.GenerateAggregateKey([]crypto.CurvePoint{pub1, pub2, pub3, pub4})

	g := game.NewGame(*context, []*big.Int{priv1, priv2, priv3, priv4}, 4)

	fmt.Printf("[Game]\ntotalPlayers: %d\ncurrentScore: %d\n---\n", g.TotalPlayers, g.Score)

	g.PlayRound(*context)
}
