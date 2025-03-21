package utils

import (
	"github.com/google/uuid"
)

func ToIntPointer(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

func ToUUIDPointer(value string) (*uuid.UUID, error) {
	if value == "" {
		return nil, nil
	}
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}
	return &parsedUUID, nil
}
