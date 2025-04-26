package game

import (
	"clyde1811/dmp/cardset"
	"clyde1811/dmp/player"
	"fmt"
	"os"
)

type Game struct {
	Players      []player.Player
	Cardset      *cardset.CardSet
	CurrentScore int
	TotalPlayers int
}

func NewGame(numPlayers int) *Game {
	// Check validity of player number
	if numPlayers < 2 || numPlayers > 5 {
		fmt.Println("Number of player is invalid. Exit.")
		os.Exit(1)
	}

	// Generate Players
	players := make([]player.Player, numPlayers)

	for i := 0; i < numPlayers; i++ {
		newplayer, err := player.GeneratePlayer(i)
		if err != nil {
			fmt.Println("Failed to generate player. Exit.")
			os.Exit(1)
		}

		players[i] = newplayer
		fmt.Printf("Player %d: \n Private Key %s \n Public Key %s\n", i, players[i].GetPrivateKey(), players[i].GetPublicKey())
	}

	// Create game Cardset
	cs := cardset.CreateCardSet(48763)

	return &Game{
		Players:      players,
		Cardset:      cs,
		CurrentScore: 0,
		TotalPlayers: numPlayers,
	}
}

func (g *Game) ShuffleEncrypt() {
	/* TODO: Each player needs to participate the process of shuffling and encryption by calling ShuffleEncrypt
	Here, we just shuffle.
	*/

	// Shuffle
	g.Cardset.ShuffleCardSet()

	// Encrypt
	// Skip here
}

func (g *Game) DealDecrypt(id int) {
	/* TODO: The card on the top of Cardset is first decrypted and posted back by all the others
	   Finally, the specified player (id) completes decryption and adds the card to their hand.
	Here, we just draw a card to the player.
	*/

	card := g.Cardset.Draw()

	g.Players[id].Hand = append(g.Players[id].Hand, card)
}

func (g *Game) dealCards() {
	for i := 0; i < g.TotalPlayers; i++ {
		for j := 0; j < 5; j++ {
			g.DealDecrypt(i)
		}
	}
}

func (g *Game) showCardset() {
	fmt.Printf("[Game] Cardset:\n")
	for i := g.Cardset.Index; i < len(g.Cardset.Card); i++ {
		c := g.Cardset.Card[i]
		fmt.Printf("%d, %s\n", i, c.ToCardString())
	}

	fmt.Println("---")
}

func (g *Game) PlayRound() {
	// Shuffle and encrypt by each player
	for i := 0; i < g.TotalPlayers; i++ {
		g.ShuffleEncrypt()
	}

	// Deal 5 cards for each player
	g.dealCards()

	g.showCardset()
	for i := 0; i < g.TotalPlayers; i++ {
		p := g.Players[i]
		p.ShowHand()
	}
}
