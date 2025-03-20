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

func ToUUIDPointer(value string) *uuid.UUID {
	// TODO: rever se é melhor retornar nil ou um erro
	if value == "" {
		return nil
	}
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return nil
	}
	return &parsedUUID
}
