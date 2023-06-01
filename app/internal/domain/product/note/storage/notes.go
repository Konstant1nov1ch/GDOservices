package storage

import (
	"GDOservice/internal/domain/product/dao"
	"GDOservice/internal/domain/product/note/model"
	"GDOservice/pkg/client/postgresql"
	db "GDOservice/pkg/client/postgresql/model"
	"GDOservice/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type NoteStorage interface {
	GetNotesByTableID(ctx context.Context, tableID int) ([]model.Note, error)
	// ToDo Другие методы для работы с table
}

type PostgreSQLNoteStorage struct {
	queryBuilder sq.StatementBuilderType
	client       dao.PostgreSQLClient
	logger       *logging.Logger
}

func NewPostgreSQLNoteStorage(client dao.PostgreSQLClient, logger *logging.Logger) *PostgreSQLNoteStorage {
	return &PostgreSQLNoteStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
	}
}

func (s *PostgreSQLNoteStorage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *PostgreSQLNoteStorage) GetNotesByTableID(ctx context.Context, tableID int) ([]model.Note, error) {
	query := s.queryBuilder.Select("id").
		Column("table_id").
		Column("category_id").
		Column("deadline").
		Column("title").
		Column("description").
		From(dao.Scheme + "." + dao.Table_note).
		Where(sq.Eq{"table_id": tableID})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, dao.Table_note, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	list := make([]model.Note, 0)

	for rows.Next() {
		p := model.Note{}
		if err = rows.Scan(
			&p.Id, &p.TableId, &p.CategoryId, &p.Deadline, &p.Title, &p.Description,
		); err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}
