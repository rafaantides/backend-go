package utils

import (
	"backend-go/internal/api/errs"
	"reflect"

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

func ValidateUUIDs(strings []string) error {
	for _, str := range strings {
		if _, err := uuid.Parse(str); err != nil {
			return err
		}
	}
	return nil
}

func ValidateUUIDArrayFields(s any) *string {
	val := reflect.ValueOf(s).Elem()

	// Percorre todos os campos da struct
	for i := range val.NumField() {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Verifica se o campo é um ponteiro para slice de strings
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Slice {
			slice := fieldValue.Elem()

			// Valida apenas se o slice não estiver vazio
			if slice.Len() > 0 {
				sliceValue := slice.Interface().([]string)
				if err := ValidateUUIDs(sliceValue); err != nil {
					detail := errs.InvalidUUID(field.Tag.Get("form"), err.Error())
					return &detail
				}
			}
		}
	}

	return nil
}
