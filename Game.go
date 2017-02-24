package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// StartNewGame initalizes the players, the deck, and deals cards to each
// player.
func StartNewGame(name *string, pScore, AIScore *int) (err error) {
	p1 := &Player{*name, Hand{}}
	p2 := &Player{"AI", Hand{}}
	turn := p1
	RummyDeck := InitializeDeck()
	RummyStack := RummyDeck.InitializeStack()
	RummyDeck.Deal(p1, p2)
	knock, draw := false, false

	// While Knock is true, keep the players in a loop that handle turns.
	for {
		if turn == p1 {
			PlayerActions(p1, &RummyDeck, &RummyStack, &knock, &draw)
			if knock || draw {
				break
			}
			turn = p2
		} else {
			AIActions(p2, &RummyDeck, &RummyStack, &knock, &draw)
			if knock || draw {
				break
			}
			turn = p1
		}
	}

	if draw {
		fmt.Printf("The game was a draw!\n---\n%s score: %d AI score: %d \n Play again? (Y/N)", *name, *pScore, *AIScore)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		response = strings.ToUpper(strings.TrimSpace(response))

		if response == "Y" {
			StartNewGame(name, pScore, AIScore)
		}
		return err
	}

	if turn == p1 {
		*pScore = CalculateScore(&p1.Hand, &p2.Hand)
	} else {
		*AIScore = CalculateScore(&p2.Hand, &p1.Hand)
	}

	fmt.Printf("%s score: %d AI score: %d \nPlay again? (Y/N)", *name, *pScore, *AIScore)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	response = strings.ToUpper(strings.TrimSpace(response))

	if response == "Y" {
		StartNewGame(name, pScore, AIScore)
	}
	fmt.Printf("Goodbye!\nFinal scores:\n%s: %d\nAI: %d", *name, *pScore, *AIScore)
	return err

}

// PlayerActions - describes what the player is going to do.
func PlayerActions(p *Player, deck *Deck, stack *Stack, knock *bool, draw *bool) {
	reader := bufio.NewReader(os.Stdin)

TURN_ACTIONS:
	for {
		if len(*deck) == 0 {
			*draw = true
			break TURN_ACTIONS
		}
		fmt.Printf("\n---\nCard on stack: %s \nYour hand: %s \n", stack.PeekAtStack(), p.PrettyPrintHand())
		fmt.Printf("\nWhat would you like to do, %s?\n 1. DRAW CARD FROM DECK\n 2. PICKUP CARD FROM STACK\n 3. CHECK MELDS IN HAND\n 4. CHECK POINTS IN HAND\n", p.name)
		response, err := reader.ReadString('\n')
		response = strings.TrimRight(response, "\n")
		if err != nil {
			fmt.Println("Unrecognized command.")
			continue
		}
		response = strings.ToUpper(strings.TrimSpace(response))
		switch response {
		case "1", "DRAW CARD":
			p.Hand.DrawCard(deck)
			break TURN_ACTIONS
		case "2", "PICKUP CARD FROM STACK":
			p.Hand.DrawCard(stack)
			break TURN_ACTIONS
		case "3", "CHECK MELDS":
			melds := p.Hand.CheckMelds()
			fmt.Printf("\n%s", melds.PrettyPrintMelds())
		case "4", "CHECK POINTS":
			// Check the total of points in your hand, values not melded
			total := p.Hand.CheckTotal()
			if total <= 10 {
				fmt.Printf("\nYour hand total is: %d. Will you knock? (Y/N) ", total)
				response, err := reader.ReadString('\n')
				response = strings.ToUpper(strings.TrimRight(response, "\n"))

				if err != nil {
					fmt.Println("Something went wrong...")
				}

				if response == "Y" {
					*knock = true
					break TURN_ACTIONS
				}
			}

			fmt.Printf("\nYou hand total is: %d", total)
		default:
			fmt.Printf("\nUnrecognized command. Try again.\n\n")
		}
	}

	if *knock {
		return
	}

	if *draw {
		return
	}

	for {
		fmt.Printf("\nYou must now discard a card from your hand.\n\nCard on stack: %s \nYour hand: %s \n", stack.PeekAtStack(), p.PrettyPrintHand())

		discard, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Something went wrong. Try again.")
			continue
		}

		discard = strings.ToUpper(strings.TrimRight(discard, "\n"))
		card, err := GetCardFromPrettyPrint(discard)
		if err != nil {
			fmt.Println("Something went wrong. Try again.")
			continue
		}
		p.Hand.DiscardCard(card, stack)
		break
	}
	return
}

// CalculateScore - gets the score from the last round.
func CalculateScore(h1, h2 *Hand) (score int) {
	return h2.CheckTotal() - h1.CheckTotal()
}

func main() {
	pScore := 0
	AIScore := 0

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your name:")
	name, err := reader.ReadString('\n')
	name = strings.TrimRight(name, "\n")
	fmt.Printf("\n")

	if err != nil {
		fmt.Printf("An error occurred with the name")
		os.Exit(0)
	}

	fmt.Printf("Welcome %s! Lets play a game of Gin Rummy!\n", name)
	fmt.Printf("\n")
	StartNewGame(&name, &pScore, &AIScore)
}
