package utils

import (
	"github.com/google/uuid"
)

// TODO: rever a utilização dessa funçao
func ToUUIDPointer(str string) (*uuid.UUID, error) {
	if str == "" {
		return nil, nil
	}
	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return nil, err
	}
	return &parsedUUID, nil
}
