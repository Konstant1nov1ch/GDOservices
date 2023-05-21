package handler

import (
	"GDOservice/internal/domain/product/user/storage"
	db "GDOservice/pkg/client/postgresql/model"
	"encoding/json"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(userStorage storage.UserStorage) http.HandlerFunc {
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

		response, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
