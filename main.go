package main

import (
	"clyde1811/dmp/cardset"
	"clyde1811/dmp/player"
	"fmt"
)

func main() {
	fmt.Println("Decentralized Mental Poker v1.0");
	fmt.Println("We have four players in default.");

	var players [4]player.Player;

	// Create card set
	cardset := cardset.CreateCardSet(48763)

	// Shuffle card
	cardset.ShuffleCardSet()

	for i := 0; i < 52; i++ {
		// Draw card
		c := cardset.Draw()

		for j := 0; j < 4; j++ {
			if i % 4 == j {
				continue
			}
			players[j].EstablishCard(c)
		}

		fmt.Printf("Player %d draw %s\n", (i)%4+1, c.ToCardString())
	}
}