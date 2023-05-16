package dao

import (
	"GDOservice/internal/domain/product/table/model"
	"GDOservice/pkg/client/postgresql"
	db "GDOservice/pkg/client/postgresql/model"
	"GDOservice/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type TableStorage struct {
	queryBuilder sq.StatementBuilderType
	client       PostgreSQLClient
	logger       *logging.Logger
}

func NewTableStorage(client PostgreSQLClient, logger *logging.Logger) TableStorage {
	return TableStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
	}
}

func (s *TableStorage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *TableStorage) AllTablesByUserID(ctx context.Context, userID int) ([]model.Table, error) {
	query := s.queryBuilder.Select("id").
		Column("capacity").
		From(scheme + "." + table_table).
		Where(sq.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table_table, args)
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

	list := make([]model.Table, 0)

	for rows.Next() {
		t := model.Table{}
		if err = rows.Scan(
			&t.Id, &t.UserId, &t.Capacity,
		); err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err
		}

		list = append(list, t)
	}

	return list, nil
}
