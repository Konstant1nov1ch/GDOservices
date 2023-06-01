package note

import (
	"GDOservice/internal/domain/auth"
	"GDOservice/internal/domain/auth/cache"
	"GDOservice/internal/domain/product/note/storage"
	"encoding/json"
	"github.com/jackc/pgtype"
	"net/http"
	"strconv"
)

func NotesByTableID(noteStorage storage.NoteStorage, tokenCache cache.Cache) http.HandlerFunc {
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

		// Получение table_id из параметров запроса
		tableID := r.URL.Query().Get("table_id")
		if tableID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Преобразование tableID в int
		tableIDInt, err := strconv.Atoi(tableID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userUUID := &pgtype.UUID{}
		err = userUUID.Set(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Получение заметок по table_id
		notes, err := noteStorage.GetNotesByTableID(r.Context(), tableIDInt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Отправка списка заметок в формате JSON
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(notes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
