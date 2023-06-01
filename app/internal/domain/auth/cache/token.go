package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

const tokenLength = 5

func GenerateToken(uuid pgtype.UUID) (string, error) {
	s := fmt.Sprintf("%x", uuid.Bytes)

	hash := sha256.Sum256([]byte(s))
	token := hex.EncodeToString(hash[:tokenLength])
	return token, nil
}

func GetTokenFromCache(cache Cache, key string, uuid pgtype.UUID) (string, error) {
	token := cache.Get(key)
	if token != nil {
		return token.(string), nil
	}
	newToken, err := GenerateToken(uuid)
	if err != nil {
		return "", err
	}

	cache.Set(key, newToken)
	return newToken, nil
}

func SetTokenInCache(cache Cache, key, token string) {
	cache.Set(key, token)
}

func IsTokenExpired(creationTime time.Time) bool {
	expirationTime := creationTime.Add(30 * time.Minute)
	return time.Now().After(expirationTime)
}

func RefreshToken(cache Cache, key string, uuid pgtype.UUID) (string, error) {
	newToken, err := GenerateToken(uuid)
	if err != nil {
		return "", err
	}
	SetTokenInCache(cache, key, newToken)
	return newToken, nil
}
