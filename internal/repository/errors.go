package repository

import (
	"errors"

	"github.com/SamVitebsk/sharetrip-contract/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	postgresUniqueViolation     = "23505"
	postgresForeignKeyViolation = "23503"
	postgresNotNullViolation    = "23502"
	postgresCheckViolation      = "23514"
)

var (
	ErrInvalidReference      = errors.New("ссылка на несуществующую запись")
	ErrConstraintViolation   = errors.New("нарушено ограничение данных")
	ErrEmptyContractServices = errors.New("список услуг договора пуст")
	ErrInvalidContractData   = errors.New("некорректные данные договора в хранилище")
)

func mapPostgresError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.Join(service.ErrNotFound, err)
	}
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return err
	}
	switch postgresError.Code {
	case postgresUniqueViolation:
		return errors.Join(service.ErrConflict, err)
	case postgresForeignKeyViolation:
		return errors.Join(ErrInvalidReference, err)
	case postgresNotNullViolation, postgresCheckViolation:
		return errors.Join(ErrConstraintViolation, err)
	default:
		return err
	}
}
