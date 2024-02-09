package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"math/big"
	"time"
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

// Returns a randomized session ID
func GenerateSessionID() string {
    b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

// Checks if the access token stored in Redis identified by `tokenName` argument has expired
func IsAccessTokenExpired(appCtx *AppContext, tokenName string) (bool, error) {
    accessTokenTTL, err := appCtx.RedisClient.TTL(context.Background(), tokenName).Result()
    if err != nil {
        log.Printf("There was an issue retrieving the access token time to live: %v\n", err)
        return true, err
    }

    // Token has a TTL and is soon to expire
    if accessTokenTTL > 0 && accessTokenTTL < 5 * time.Minute {
        return true, nil
    }

    if accessTokenTTL == -2 {
        log.Printf("Access token key does not exist\n")
        return true, nil
    }

    if accessTokenTTL == -1 {
        log.Printf("Access token does not have an expiry set\n")
        return false, nil
    }
    
    // Token is not expired
    return false, nil
}