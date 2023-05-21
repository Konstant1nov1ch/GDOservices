package storage

import (
	"GDOservice/internal/domain/product/category/model"
	"GDOservice/internal/domain/product/dao"
	"GDOservice/pkg/client/postgresql"
	db "GDOservice/pkg/client/postgresql/model"
	"GDOservice/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type CategoryStorage struct {
	queryBuilder sq.StatementBuilderType
	client       dao.PostgreSQLClient
	logger       *logging.Logger
}

func NewCategoryStorage(client dao.PostgreSQLClient, logger *logging.Logger) CategoryStorage {
	return CategoryStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
	}
}

func (s *CategoryStorage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *CategoryStorage) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	query := s.queryBuilder.Select("id", "name", "back_color", "word_color").
		From(dao.Scheme + "." + dao.Table_category)
	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, dao.Table_category, args)
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

	list := make([]model.Category, 0)

	for rows.Next() {
		c := model.Category{}
		if err = rows.Scan(
			&c.Id, &c.Name, &c.BackColor, &c.WordColor,
		); err != nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err
		}

		list = append(list, c)
	}

	return list, nil
}
