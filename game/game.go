package game

import (
	"clyde1811/dmp/cardset"
	"clyde1811/dmp/crypto"
	"clyde1811/dmp/player"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

type Game struct {
	Players      []player.Player
	Cardset      *cardset.CardSet
	OrderReverse bool
	Score        int
	TotalPlayers int
	Winners      []player.Player
}

func (g *Game) IsEffectCard(card cardset.CardPoint) bool {
	return card.Number == "J" || card.Number == "Q" || card.Number == "K"
}

func (g *Game) CardValue(card cardset.CardPoint) int {
	switch card.Number {
	case "2", "3", "4", "5", "6", "7", "8", "9":
		// Convert the string digit to an int
		num, _ := strconv.Atoi(card.Number)
		return num
	case "A":
		return 1
	case "T":
		return 10
	case "Q":
		return -20
	}

	return 0
}

func (g *Game) CanPlayCard(card cardset.CardPoint) bool {
	switch card.Number {
	case "2", "3", "4", "5", "6", "7", "8", "9":
		// Convert the string digit to an int
		num, _ := strconv.Atoi(card.Number)
		return g.Score+num <= 99
	case "A":
		return g.Score+1 <= 99
	case "T":
		return g.Score+10 <= 99
	default:
		return true
	}
}

func NewGame(context crypto.CryptoContext, privKeys []*big.Int, numPlayers int) *Game {
	// Check validity of player number
	if numPlayers < 2 || numPlayers > 5 {
		fmt.Println(os.Stderr, "invalid number of players")
		os.Exit(1)
	}

	// Generate Players
	players := make([]player.Player, numPlayers)
	for i := range players {
		p := player.GeneratePlayer(context, privKeys[i], i)

		players[i] = p
		fmt.Printf("Player %d: \n Private Key %s \n Public Key (%s, %s)\n", i, players[i].GetPrivateKey().String(), players[i].GetPublicKey().X.String(), players[i].GetPublicKey().Y.String())
	}

	// Create game Cardset
	cs := cardset.CreateCardSet(context, 77149)

	return &Game{
		Players:      players,
		Cardset:      cs,
		Score:        0,
		OrderReverse: false,
		TotalPlayers: numPlayers,
		Winners:      []player.Player{},
	}
}

func (g *Game) ShuffleEncrypt() error {
	/* TODO: Each player needs to participate the process of shuffling and encryption by calling ShuffleEncrypt
	Here, we just shuffle.
	*/

	// Shuffle
	g.Cardset.ShuffleCardSet()

	// Encrypt
	// Skip here

	return nil
}

func (g *Game) DealDecrypt(context crypto.CryptoContext, id int) error {
	/* TODO: The card on the top of Cardset is first decrypted and posted back by all the others
	   Finally, the specified player (id) completes decryption and adds the card to their hand.
	Here, we just draw a card to the player.
	*/

	card, err := g.Cardset.Draw()

	if err != nil {
		return err
	}

	for _, player := range g.Players {
		if player.Id == id {
			continue
		}
		Cb := context.DecryptCard(card.Ca, card.Cb, player.PrivateKey)
		card.Cb = Cb
	}

	searchableCard := g.Players[id].DecryptCard(context, card)

	cardPoint := g.Cardset.FindCardByPoint(searchableCard.Cb)

	g.Players[id].Hand = append(g.Players[id].Hand, cardPoint)

	return nil
}

func (g *Game) DealCards(context crypto.CryptoContext) {
	for i := 0; i < g.TotalPlayers; i++ {
		for j := 0; j < 5; j++ {
			err := g.DealDecrypt(context, i)
			if err != nil {
				if errors.Is(err, cardset.ErrCardsetEmpty) {
					fmt.Println(err, ", therefore stop dealing card")
					return
				} else {
					fmt.Println(os.Stderr, "failed to deal and decrypt card with error: ", err)
					os.Exit(1)
				}
			}
		}
	}
}

func (g *Game) ShowCardset() {
	fmt.Printf("[Game] Cardset:\n")
	for i, c := range g.Cardset.CardPoint {
		fmt.Printf("%d, %s\n", i, c.ToCardString())
	}

	fmt.Println("---")
}

func (g *Game) ApplyCard(card cardset.CardPoint, playerID int) {
	switch card.Number {
	case "2", "3", "4", "5", "6", "7", "8", "9":
		num, _ := strconv.Atoi(card.Number)
		g.Score += num
		fmt.Printf("Player %d plays %s: Score + %d = %d\n", playerID, card.ToCardShrtnString(), num, g.Score)
	case "A":
		g.Score += 1
		fmt.Printf("Player %d plays %s: Score + 1 = %d\n", playerID, card.ToCardShrtnString(), g.Score)
	case "T":
		g.Score += 10
		fmt.Printf("Player %d plays %s: Score + 10 = %d\n", playerID, card.ToCardShrtnString(), g.Score)
	case "J":
		g.OrderReverse = !g.OrderReverse
		fmt.Printf("Player %d plays %s: Turn order reversed\n", playerID, card.ToCardShrtnString())
	case "Q":
		g.Score -= 20
		fmt.Printf("Player %d plays %s: Score - 20 = %d\n", playerID, card.ToCardShrtnString(), g.Score)
	case "K":
		g.Score = 0
		fmt.Printf("Player %d plays %s: Score = 0\n", playerID, card.ToCardShrtnString())
	}
}

func (g *Game) Eliminate(player *player.Player) {
	player.Active = false
}

func (g *Game) PlayRound(context crypto.CryptoContext) {
	// Shuffle and encrypt by each player
	for i := 0; i < g.TotalPlayers; i++ {
		err := g.ShuffleEncrypt()
		if err != nil {
			os.Exit(1)
		}
	}

	// Deal 5 cards for each player
	g.DealCards(context)

	// Print
	g.ShowCardset()
	for i := 0; i < g.TotalPlayers; i++ {
		p := g.Players[i]
		p.ShowHand()
	}

	// Start round
	activePlayers := g.TotalPlayers
	turn := 0
	for activePlayers > 0 {
		fmt.Printf("(turn %d) ", turn)
		// Find current player
		index := turn % g.TotalPlayers
		if g.OrderReverse == true {
			index = g.TotalPlayers - 1 - index
		}
		currentPlayer := &g.Players[index]

		if !currentPlayer.Active { // Check active
			turn++
			continue
		}

		//fmt.Printf("\nPlayer %d's turn (Score: %d)\n", currentPlayer.Id, g.Score)

		choosed, cardIdx := currentPlayer.ChooseCard(g)

		//fmt.Printf("\nChosen cardIdx: %d, %s\n", cardIdx, choosed.ToCardString())

		// Eliminate the player if it cannot play any card
		if cardIdx == -1 {
			g.Eliminate(currentPlayer)
			activePlayers--
			fmt.Printf("Eliminating player %d\n", currentPlayer.Id)
			turn++
			continue
		}

		g.ApplyCard(choosed, currentPlayer.Id)

		// Draw a card if deck is not empty
		err := g.DealDecrypt(context, currentPlayer.Id)
		if err != nil {
			if errors.Is(err, cardset.ErrCardsetEmpty) {
				// do nothing
			} else {
				fmt.Println(os.Stderr, "failed to deal and decrypt card with error: ", err)
				os.Exit(1)
			}
		}

		// Check if player wins
		if len(currentPlayer.Hand) == 0 && g.Score < 99 {
			currentPlayer.Active = false
			g.Winners = append(g.Winners, *currentPlayer)
			activePlayers--
			fmt.Printf("Player %d wins the round!\n", currentPlayer.Id)
		}

		// Determine next player
		if choosed.Number != "J" {
			turn++
		}
	}
}
