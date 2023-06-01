package storage

import (
	"GDOservice/internal/domain/auth/cache"
	noteStorage "GDOservice/internal/domain/product/note/storage"
	handler "GDOservice/internal/domain/product/table"
	"GDOservice/internal/domain/product/table/model"
	tableStorage "GDOservice/internal/domain/product/table/storage"
	userStorage "GDOservice/internal/domain/product/user/storage"
	"GDOservice/pkg/client/postgresql"
	"GDOservice/pkg/logging"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNoteTablesByUser(t *testing.T) {
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
	userStorage := userStorage.NewPostgreSQLUserStorage(client, &logger)

	tokenCache, err := cache.NewCache(10, 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	email := "1234"
	password := "1234"

	user, err := userStorage.AuthenticateUser(ctx, email, password)

	actualToken, err := cache.GenerateToken(user.Id)
	if err != nil {
		logger.Println("Failed to create token")
	}
	userSID := fmt.Sprintf("%x-%x-%x-%x-%x", user.Id.Bytes[0:4], user.Id.Bytes[4:6], user.Id.Bytes[6:8], user.Id.Bytes[8:10], user.Id.Bytes[10:16])

	cache.SetTokenInCache(tokenCache, actualToken, userSID)

	// Создание экземпляра хранилища таблиц
	tableStorage := *tableStorage.NewPostgreSQLTableStorage(client, &logger)

	// Создание экземпляра хранилища заметок
	noteStorage := noteStorage.NewNoteStorage(client, &logger)

	// Создание экземпляра хендлера
	handler := handler.TablesByUser(&tableStorage, tokenCache)

	// Создание фейкового HTTP-запроса
	req, err := http.NewRequest("GET", "/tables", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Установка заголовка Authorization с токеном
	req.Header.Set("Authorization", actualToken)

	// Создание фейкового HTTP-ответа
	recorder := httptest.NewRecorder()

	// Выполнение HTTP-запроса с помощью хендлера
	handler.ServeHTTP(recorder, req)

	// Проверка статусного кода HTTP-ответа
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, recorder.Code)
	}

	// Проверка Content-Type заголовка
	expectedContentType := "application/json"
	actualContentType := recorder.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Expected Content-Type %s but got %s", expectedContentType, actualContentType)
	}

	// Проверка содержимого тела ответа
	expectedResponse := []model.Table{
		{Id: 1, Capacity: 2},
	}
	var actualResponse []model.Table
	err = json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if len(actualResponse) != len(expectedResponse) {
		t.Errorf("Expected %d tables but got %d", len(expectedResponse), len(actualResponse))
	}
	logger.Println(actualResponse)

	// Получение списка заметок для каждой таблицы
	for _, table := range actualResponse {
		notes, err := noteStorage.GetNotesByTableID(ctx, table.Id)
		if err != nil {
			t.Errorf("Failed to get notes for table %d: %v", table.Id, err)
		}
		logger.Println(notes)
	}
}
