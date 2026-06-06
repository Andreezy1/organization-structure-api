package repository

import (
	"errors"
	"org_struct_api/internal/models"

	"gorm.io/gorm"
)

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ErrNotFound
	}
	return err
}
