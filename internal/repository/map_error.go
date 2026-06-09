package repository

import (
	"errors"
	"org_struct_api/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const pgErrCodeUniqueViolation = "23505"

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeUniqueViolation {
		return models.ErrConflict
	}
	return err
}
