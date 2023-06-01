package storage

import (
	"context"

	"GDOservice/internal/domain/product/note/model"
)

type MockNoteStorage struct {
	GetNotesByTableIDFunc func(ctx context.Context, tableID int) ([]model.Note, error)
}

func (m *MockNoteStorage) GetNotesByTableID(ctx context.Context, tableID int) ([]model.Note, error) {
	return m.GetNotesByTableIDFunc(ctx, tableID)
}
