package auth

import (
	"GDOservice/internal/domain/auth/cache"
	"fmt"
)

func IsValidToken(token string, tokenCache cache.Cache) bool {
	// Проверка токена на валидность, например, в кеше или базе данных
	return tokenCache.Get(token) != nil
}

func GetUserIDFromToken(token string, tokenCache cache.Cache) (string, error) {
	userID := tokenCache.Get(token)
	if userID == nil {
		return "", fmt.Errorf("invalid token")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("invalid user_id type")
	}

	return userIDStr, nil
}
