package random

import (
	"math/rand"
	"time"
)

// Constant letters for randomize
const (
	Letters     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	DigiLetters = "1234567890"
)

// GenerateRandomString return random string by n
func GenerateRandomString(n int) string {
	letters := []rune(Letters)
	rand.Seed(time.Now().UTC().UnixNano())
	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

// GenerateRandomDigitString return random digit string by n
func GenerateRandomDigitString(n int) string {
	letters := []rune(DigiLetters)
	rand.Seed(time.Now().UTC().UnixNano())
	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}
