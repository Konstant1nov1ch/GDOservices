package storage

import (
	"GDOservice/internal/domain/product/dao"
	"GDOservice/internal/domain/product/user/model"
	"GDOservice/pkg/client/postgresql"
	db "GDOservice/pkg/client/postgresql/model"
	"GDOservice/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type UserStorage interface {
	AuthenticateUser(ctx context.Context, email, password string) (*model.User, error)
	// ToDo Другие методы для работы с пользователями
}

type PostgreSQLUserStorage struct {
	queryBuilder sq.StatementBuilderType
	client       dao.PostgreSQLClient
	logger       *logging.Logger
}

func NewPostgreSQLUserStorage(client dao.PostgreSQLClient, logger *logging.Logger) *PostgreSQLUserStorage {
	return &PostgreSQLUserStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
	}
}

func (s *PostgreSQLUserStorage) queryLogger(sql string, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

// ToDo добавить хэш
func (s *PostgreSQLUserStorage) AuthenticateUser(ctx context.Context, email, password string) (*model.User, error) {
	query := s.queryBuilder.Select("email", "pwd", "name", "id", "payment_status").
		From(dao.Scheme + "." + dao.Table_user).
		Where(sq.Eq{"email": email, "pwd": password})
	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, dao.Table_user, args)
	if err != nil {
		err = db.ErrCreateQuery(err)
		logger.Error(err)
		return nil, err
	}

	logger.Trace("do query")
	row := s.client.QueryRow(ctx, sql, args...)

	user := model.User{}
	if err := row.Scan(&user.Email, &user.Pwd, &user.Name, &user.Id, &user.PaymentStatus); err != nil {
		if err == nil {
			err = db.ErrScan(postgresql.ParsePgError(err))
			logger.Error(err)
			return nil, err // Пользователь не найден
		}
		err = db.ErrScan(err)
		logger.Error(err)
		return nil, err
	}

	return &user, nil
}
