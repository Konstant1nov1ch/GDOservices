package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	handlerStorage "GDOservice/internal/domain/product/note"
	"GDOservice/internal/domain/product/note/model"

	"github.com/stretchr/testify/assert"
)

type Cache interface {
	Set(key string, value string)
	Get(key string) interface{}
}

type NoteStorage interface {
	GetNotesByTableID(ctx context.Context, tableID int) ([]model.Note, error)
}

type mockCache struct {
	data map[string]interface{}
}

func (m *mockCache) Set(key string, value string) {
	m.data[key] = value
}

func (m *mockCache) Get(key string) interface{} {
	value, ok := m.data[key]
	if !ok {
		return nil
	}
	return value
}

type mockNoteStorage struct {
	GetNotesByTableIDFunc func(ctx context.Context, tableID int) ([]model.Note, error)
}

func (m *mockNoteStorage) GetNotesByTableID(ctx context.Context, tableID int) ([]model.Note, error) {
	return m.GetNotesByTableIDFunc(ctx, tableID)
}

func TestNotesByTableID(t *testing.T) {
	//ToDo fix EOF exp
	// Create a mock NoteStorage
	mockNoteStorage := &mockNoteStorage{
		GetNotesByTableIDFunc: func(ctx context.Context, tableID int) ([]model.Note, error) {
			// Return fake notes
			return []model.Note{
				{Id: 1, TableId: tableID, Title: "Note 1"},
				{Id: 2, TableId: tableID, Title: "Note 2"},
			}, nil
		},
	}

	// Create a mock TokenCache
	mockTokenCache := &mockCache{
		data: make(map[string]interface{}),
	}

	// Set a fake token in the cache
	mockTokenCache.Set("1fb5327de8", "1234")

	// Create a test server and handler
	server := httptest.NewServer(http.HandlerFunc(handlerStorage.NotesByTableID(mockNoteStorage, mockTokenCache)))
	defer server.Close()

	// Create a GET request to the test server
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?table_id=1", server.URL), nil)
	assert.NoError(t, err)

	// Set the Authorization header
	req.Header.Set("Authorization", "1fb5327de8")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the JSON response
	var notes []model.Note
	err = json.NewDecoder(resp.Body).Decode(&notes)
	assert.NoError(t, err)

	// Check the received notes
	assert.Len(t, notes, 2)
	assert.Equal(t, 1, notes[0].Id)
	assert.Equal(t, 1, notes[0].TableId)
	assert.Equal(t, "Note 1", notes[0].Title)
	assert.Equal(t, 2, notes[1].Id)
	assert.Equal(t, 1, notes[1].TableId)
	assert.Equal(t, "Note 2", notes[1].Title)
}
