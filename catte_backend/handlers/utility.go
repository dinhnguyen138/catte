package handlers

import (
	"math/rand"
	"time"

	"github.com/dinhnguyen138/catte/catte_backend/constants"
	"github.com/dinhnguyen138/catte/catte_backend/models"
)

func deal() (deck models.Deck) {
	// Loop over each type and suit appending to the deck
	for i := 0; i < len(constants.Types); i++ {
		for n := 0; n < len(constants.Suits); n++ {
			deck = append(deck, constants.Types[i]+constants.Suits[n])
		}
	}
	return
}

// Shuffle the deck
func shuffle(d models.Deck) models.Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
	return d
}

func larger(leftCard string, rightCard string) bool {
	if rightCard == "" {
		return true
	}
	leftValue := leftCard[:len(leftCard)-1]
	leftSuit := leftCard[len(leftCard)-1:]
	rightValue := rightCard[:len(rightCard)-1]
	rightSuit := rightCard[len(rightCard)-1:]
	if leftSuit == rightSuit {
		return constants.CardOrder[leftValue] > constants.CardOrder[rightValue]
	}
	return false
}
