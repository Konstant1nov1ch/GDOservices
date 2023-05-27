package storage_test

import (
	"GDOservice/internal/domain/product/user/storage"
	"GDOservice/pkg/client/postgresql"
	"GDOservice/pkg/logging"
	"context"
	"testing"
	"time"
)

func TestAuthenticateUser(t *testing.T) {
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

	user, err := storage.AuthenticateUser(ctx, email, password)
	if err != nil {
		t.Fatalf("Failed to authenticate user: %v", err)
	}

	if user == nil {
		t.Fatal("User not found or authentication failed")
	}
}
