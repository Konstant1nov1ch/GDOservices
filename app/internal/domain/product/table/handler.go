package handler

import (
	"GDOservice/internal/domain/auth/cache"
	tableStorage "GDOservice/internal/domain/product/table/storage"
	userModel "GDOservice/internal/domain/product/user/model"
	userStorage "GDOservice/internal/domain/product/user/storage"
	"encoding/json"
	"net/http"
)

func TablesByUser(userStorage *userStorage.UserStorage, tokenCache *cache.Cache, tblStorage *tableStorage.TableStorage) http.HandlerFunc {
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

		// Получение пользователя по токену
		user, err := GetUserFromToken(token, userStorage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Обновление токена пользователя в кэше
		err = RefreshUserToken(token, user, tokenCache)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Получение таблиц пользователя
		tables, err := tblStorage.AllTablesByUserID(r.Context(), user.Id)
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

func GetUserFromToken(tokenCache *cache.Cache, token string) (*userModel.User, error) {
	// Получение пароля пользователя из кеша по токену
	pwd, err := cache.GetPwdFromCache(tokenCache, token)
	if err != nil {
		return nil, err
	}
	//ToDo фикс срочно
	user, err := userStorage.GetUserByPassword(pwd)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func RefreshUserToken(token string, user *userModel.User, tokenCache *cache.Cache) error {
	// Обновление токена пользователя в кэше
	_, err := cache.RefreshToken(tokenCache, token, user.Id)
	return err
}
