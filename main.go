package main

import (
	"clyde1811/dmp/cardset"
	"clyde1811/dmp/player"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Decentralized Mental Poker v1.0")
	fmt.Println("We have four players in default.")

	var players [4]player.Player

	for i := 0; i < 4; i++ {
		player, err := player.GeneratePlayer()
		if err != nil {
			fmt.Println("Failed to generate player. Exit.")
			os.Exit(1)
		}
		players[i] = player
		fmt.Printf("%d: \n Private Key %s \n Public Key %s\n", i, players[i].GetPublicKey(), players[i].GetPrivateKey())
	}

	// Create card set
	cs := cardset.CreateCardSet(48763)

	// Shuffle card
	cs.ShuffleCardSet()

	for i := 0; i < 52; i++ {
		// Draw card
		c := cs.Draw()

		for j := 0; j < 4; j++ {
			if i%4 == j {
				continue
			}
			players[j].EstablishCard(c)
		}

		cipherText, err := player.EncryptCard([]byte(c.ToCardString()), players[0].PublicKey)

		if err != nil {
			fmt.Errorf("Failed to encrypt card.")
		}

		cipherText, err = player.EncryptCard(cipherText, players[1].PublicKey)

		if err != nil {
			fmt.Errorf("Failed to encrypt card.")
		}

		plaintext, err := player.DecryptCard(cipherText, players[1].PrivateKey)

		if err != nil {
			fmt.Errorf("Failed to decrypt card.")
		}

		plaintext, err = player.DecryptCard(plaintext, players[0].PrivateKey)

		if err != nil {
			fmt.Errorf("Failed to decrypt card.")
		}

		cardString := string(plaintext)
		newCard, err := cardset.ParseCardByString(cardString)

		if err != nil {
			fmt.Errorf("Failed to parse card.")
		}

		fmt.Println(newCard.ToCardString())
	}
}
