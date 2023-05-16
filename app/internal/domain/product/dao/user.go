package dao

import (
	"GDOservice/internal/domain/product/user/model"
	"GDOservice/pkg/client/postgresql"
	db "GDOservice/pkg/client/postgresql/model"
	"GDOservice/pkg/logging"
	"context"
	sq "github.com/Masterminds/squirrel"
)

type UserStorage struct {
	queryBuilder sq.StatementBuilderType
	client       PostgreSQLClient
	logger       *logging.Logger
}

func NewUserStorage(client PostgreSQLClient, logger *logging.Logger) UserStorage {
	return UserStorage{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
		logger:       logger,
	}
}

func (s *UserStorage) queryLogger(sql, table string, args []interface{}) *logging.Logger {
	return s.logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}

func (s *UserStorage) AuthenticateUser(ctx context.Context, email, password string) (*model.User, error) {
	query := s.queryBuilder.Select("email", "pwd", "name", "id", "payment_status").
		From(scheme + "." + table_user).
		Where(sq.Eq{"email": email, "pwd": password})
	sql, args, err := query.ToSql()
	logger := s.queryLogger(sql, table_user, args)
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
