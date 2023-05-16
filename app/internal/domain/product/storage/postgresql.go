package storage

import (
	"GDOservice/pkg/logging"
	sq "github.com/Masterminds/squirrel"
)

type ProductStorage struct {
	queryBuilder sq.StatementBuilderType
	client       PostgreSQLClient
	logger       *logging.Logger
}
