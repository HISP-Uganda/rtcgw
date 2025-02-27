package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
const allowedCharacters = "0123456789" + alphabet
const codeSize = 11

// GenerateUID return a Unique ID for our resources
func GenerateUID() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source) // Creates a new instance of rand.Rand, safe for concurrent use

	numberOfCodePoints := len(allowedCharacters)

	var s strings.Builder
	s.Grow(codeSize) // Pre-allocate memory to improve performance

	// Ensure the first character is an uppercase letter from the alphabet
	s.WriteByte(allowedCharacters[r.Intn(26)] - 32) // Convert to uppercase

	// Generate the rest of the UID
	for i := 1; i < codeSize; i++ {
		s.WriteByte(allowedCharacters[r.Intn(numberOfCodePoints)])
	}

	return s.String()
}
