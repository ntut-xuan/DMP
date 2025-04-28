package main

import (
	"clyde1811/dmp/game"
	"fmt"
)

func main() {
	fmt.Println("Decentralized Mental Poker v1.0")
	fmt.Println("We have four players in default.")

	g := game.NewGame(4)

	fmt.Printf("[Game]\ntotalPlayers: %d\ncurrentScore: %d\n---\n", g.TotalPlayers, g.Score)

	g.PlayRound()
}
