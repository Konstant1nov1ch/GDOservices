package storage

import (
	"GDOservice/internal/domain/product/dao"
	"GDOservice/internal/domain/product/table/model"
	"GDOservice/pkg/client/postgresql"
	db "GDOservice/pkg/client/postgresql/model"
	"GDOservice/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

type TableStorage struct {
	queryBuilder sq.StatementBuilderType
	client       dao.PostgreSQLClient
	logger       *logging.Logger
}

func NewTableStorage(client dao.PostgreSQLClient, logger *logging.Logger) TableStorage {
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

func (s *TableStorage) AllTablesByUserID(ctx context.Context, userID pgtype.UUID) ([]model.Table, error) {
	query := s.queryBuilder.Select("id").
		Column("capacity").
		From(dao.Scheme + "." + dao.Table_table).
		Where(sq.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, dao.Table_table, args)
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
			&t.Id, &t.Capacity,
		); err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err
		}

		list = append(list, t)
	}

	return list, nil
}
