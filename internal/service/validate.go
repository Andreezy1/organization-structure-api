package service

import (
	"fmt"
	"org_struct_api/internal/models"
	"strings"
	"unicode/utf8"
)

const maxStringLen = 200

func validate(field, value string) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", fmt.Errorf("%w: %s is required", models.ErrValidation, field)
	}
	if utf8.RuneCountInString(v) > maxStringLen {
		return "", fmt.Errorf("%w: %s must not exceed %d characters", models.ErrValidation, field, maxStringLen)
	}
	return v, nil
}
