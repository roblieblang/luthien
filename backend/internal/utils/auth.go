package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
)

// Returns a high-entropy random string which will be used as a code verifier after being hashed
func GenerateCodeVerifier(length int) (string, error) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.-~"
    result := make([]byte, length)
    for i := range result {
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
        if err != nil {
            return "", err
        }
        result[i] = chars[num.Int64()]
    }
    return string(result), nil
}

// Returns a SHA256 hashed and base64 encoded string which can be used as a code verifier
func SHA256Hash(data string) string {
    hash := sha256.Sum256([]byte(data))
    return base64.RawURLEncoding.EncodeToString(hash[:])
}
