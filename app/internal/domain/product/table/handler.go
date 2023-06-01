package handler

import (
	"GDOservice/internal/domain/auth"
	"GDOservice/internal/domain/auth/cache"
	tableStorage "GDOservice/internal/domain/product/table/storage"
	"encoding/json"
	"github.com/jackc/pgtype"
	"net/http"
)

func TablesByUser(tblStorage tableStorage.TableStorage, tokenCache cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение токена из заголовка Authorization
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Проверка токена на валидность
		if !auth.IsValidToken(token, tokenCache) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Извлечение user_id из токена
		userID, err := auth.GetUserIDFromToken(token, tokenCache)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
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
