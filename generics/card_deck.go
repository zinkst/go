package main

import (
	"fmt"
	"math/rand"
	"time"
)

type PlayingCard struct {
	Suit string
	Rank string
}

func NewPlayingCard(suit string, card string) *PlayingCard {
	return &PlayingCard{Suit: suit, Rank: card}
}

func (pc *PlayingCard) String() string {
	return fmt.Sprintf("%s of %s", pc.Rank, pc.Suit)
}

type Deck[cardGeneric any] struct {
	cards []cardGeneric
}

func NewPlayingCardDeck() *Deck[*PlayingCard] {
	suits := []string{"Diamonds", "Hearts", "Clubs", "Spades"}
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	deck := &Deck[*PlayingCard]{}
	for _, suit := range suits {
		for _, rank := range ranks {
			deck.AddCard(NewPlayingCard(suit, rank))
		}
	}
	return deck
}

func (d *Deck[cardGeneric]) AddCard(card cardGeneric) {
	d.cards = append(d.cards, card)
}

func (d *Deck[cardGeneric]) RandomCard() cardGeneric {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	cardIdx := r.Intn(len(d.cards))
	return d.cards[cardIdx]
}

func main() {
	deck := NewPlayingCardDeck()

	fmt.Printf("--- drawing playing card ---\n")
	playingCard := deck.RandomCard()
	fmt.Printf("drew card: %s\n", playingCard)

	fmt.Printf("card suit: %s\n", playingCard.Suit)
	fmt.Printf("card rank: %s\n", playingCard.Rank)
}
