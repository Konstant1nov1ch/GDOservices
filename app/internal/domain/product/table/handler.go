package note

import (
	"GDOservice/internal/domain/product/user/storage"
	"net/http"
)

func TablesByUser(userStorage storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
