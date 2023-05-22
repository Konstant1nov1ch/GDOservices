package handler

import (
	"GDOservice/internal/domain/auth/cache"
	"GDOservice/internal/domain/product/user/storage"
	db "GDOservice/pkg/client/postgresql/model"
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"pwd"`
}

func LoginHandler(userStorage storage.UserStorage, tokenCache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		authRequest := AuthRequest{}
		if err := decoder.Decode(&authRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if the token exists in the cache
		token := tokenCache.Get(authRequest.Email)
		if token != nil {
			// Token exists in the cache, return it
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(token.(string)))
			return
		}

		// Token doesn't exist in the cache, authenticate user
		user, err := userStorage.AuthenticateUser(r.Context(), authRequest.Email, authRequest.Password)
		if err != nil {
			if appErr, ok := err.(*db.Error); ok {
				if appErr.Code == db.CodeNotFound {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Generate a new token
		newToken, err := cache.GenerateToken(user.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token = newToken

		// Store the token in the cache
		tokenCache.Set(authRequest.Email, token.(string))

		response, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
