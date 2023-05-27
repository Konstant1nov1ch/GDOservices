package handler

import (
	"GDOservice/internal/domain/auth/cache"
	tableStorage "GDOservice/internal/domain/product/table/storage"
	"encoding/json"
	"errors"
	"github.com/jackc/pgtype"
	"net/http"
)

func TablesByUser(tblStorage *tableStorage.TableStorage, tokenCache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение токена из заголовка Authorization
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Проверка токена на валидность
		if !IsValidToken(token, tokenCache) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Получение user_id по токену
		userID, err := GetUserIDFromToken(token, tokenCache)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userUUID := &pgtype.UUID{}
		err = userUUID.Set(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Получение таблиц пользователя
		tables, err := tblStorage.AllTablesByUserID(r.Context(), *userUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Отправка списка таблиц в формате JSON
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tables)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func IsValidToken(token string, tokenCache *cache.Cache) bool {
	// Проверка токена на валидность, например, в кеше или базе данных
	return tokenCache.Get(token) != nil
}

func GetUserIDFromToken(token string, tokenCache *cache.Cache) (string, error) {
	userID := tokenCache.Get(token)
	if userID == nil {
		return "", errors.New("invalid token")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user_id type")
	}

	return userIDStr, nil
}
