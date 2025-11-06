package authutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// KeyGenerator generates a random API key with the provided prefix.
func KeyGenerator(prefix string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error reading random bytes: %v", err)
	}
	key := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s", prefix, key), nil
}

// HashApiKeys hashes the ApiKey using SHA-256 and returns hex string.
func HashApiKeys(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
