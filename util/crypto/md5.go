package crypto

import (
	"crypto/md5"
	"fmt"
)

//GenerateMD5 return generated md5 hash string
func GenerateMD5(value string) string {
	// value = strings.ToLower(strings.TrimSpace(value))
	hash := md5.New()
	hash.Write([]byte(value))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
