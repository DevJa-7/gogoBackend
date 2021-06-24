package crypto

import (
	"crypto/sha256"
	"fmt"
)

//GenerateHash return generated sha hash string
func GenerateHash(value string) string {
	// value = strings.ToLower(strings.TrimSpace(value))
	hash := sha256.New()
	hash.Write([]byte(value))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
