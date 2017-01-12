package services

import (
	"crypto/rand"
	"math/big"
)

type UsernameGenerator interface {
	Generate() string
}

type englishAdjNounGenerator struct {
	adjectives []string
	nouns      []string
}

func NewUsernameGenerator() UsernameGenerator {
	return &englishAdjNounGenerator{
		adjectives: []string{
			"Talking", "Blue", "Purple", "Red", "Happy", "Particular", "Working", "Proud", "Silly", "Grand", "Busy",
			"Winter", "Fancy", "Active", "Resident", "Large", "Remarkable", "Special", "Heavy", "Fair", "Honest",
			"Brilliant", "Funny", "Thick", "Anxious", "Leading", "Deep", "Peculiar", "Odd", "Flying",
		},
		nouns: []string{
			"Banana", "Apple", "Chicken", "Orange", "House", "Art", "Wheels", "Person", "Internet", "Television",
			"Fact", "Media", "Movie", "Guy", "Lady", "Glasses", "Fish", "Monkey", "News", "Office", "Agency",
			"Device", "Drama", "Anxiety", "Bread", "Mom", "Dad", "Song", "Doctor", "Cookie", "Pizza", "Piano",
		},
	}
}

func (g *englishAdjNounGenerator) Generate() string {
	adjIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(g.adjectives))))
	nounIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(g.nouns))))

	return g.adjectives[adjIndex.Int64()] + g.nouns[nounIndex.Int64()]
}
