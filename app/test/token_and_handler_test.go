package storage

import (
	"GDOservice/internal/domain/auth/cache"
	"GDOservice/internal/domain/product/user"
	"GDOservice/internal/domain/product/user/storage"
	"GDOservice/pkg/client/postgresql"
	"GDOservice/pkg/logging"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoginHandler(t *testing.T) {
	// Создание конфигурации для клиента PostgreSQL
	cfg := postgresql.NewPgConfig("konstantin", "konstantin", "localhost", "5432", "todo")

	// Подключение к базе данных PostgreSQL
	ctx := context.Background()
	maxAttempts := 5
	maxDelay := 5 * time.Second
	client, err := postgresql.NewClient(ctx, maxAttempts, maxDelay, cfg)

	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer client.Close()

	logger := logging.GetLogger("info")
	storage := storage.NewUserStorage(client, &logger)

	email := "1234"
	password := "1234"

	// Создание кэша токенов
	tokenCache, err := cache.NewCache(10, 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to create token cache: %v", err)
	}

	user, err := storage.AuthenticateUser(ctx, email, password)
	if err != nil {
		t.Fatalf("Failed to authenticate user: %v", err)
	}

	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(string(body)))

	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Создание HTTP-ответа (используется httptest.ResponseRecorder)
	recorder := httptest.NewRecorder()

	// Вызов обработчика
	handler := handler.LoginHandler(storage, tokenCache)
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, recorder.Code)
	}

	expectedToken := "1fb5327de8"

	// Генерация токена
	actualToken, err := cache.GenerateToken(user.Id)
	if err != nil {
		logger.Println("Failed to create token")
	}

	if expectedToken != actualToken {
		logger.Println("Expected - 1fb5327de8, but got - " + actualToken)
	}

	token := tokenCache.Get(email)
	if token != nil {
		logger.Println("ok token from cache - " + token.(string))
	}

}
