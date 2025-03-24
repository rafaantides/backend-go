package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNoRows           = errors.New("registro não encontrado")
	ErrUnknown          = errors.New("erro desconhecido")
	ErrPermissionDenied = errors.New("permissão negada")
	ErrTimeout          = errors.New("tempo limite excedido")
)

func ErrInvalidField(field string) error {
	return fmt.Errorf("campo %s inválido", field)
}

func ErrNotFound(entity, key string) error {
	return fmt.Errorf("%s não encontrado para %s", entity, key)
}

func ErrParsingField(field string, err error) error {
	return fmt.Errorf("erro ao processar o campo %s: %w", field, err)
}

func ErrDateParsing(field string) error {
	return fmt.Errorf("erro ao processar o campo %s: use o formato YYYY-MM-DD", field)
}

func ErrFetchingData(entity, key string, err error) error {
	return fmt.Errorf("erro ao buscar %s para %s: %w", entity, key, err)
}

func ErrConflict(entity, key string) error {
	return fmt.Errorf("conflito ao processar %s com %s", entity, key)
}

func ErrExternalDependency(service string, err error) error {
	return fmt.Errorf("erro ao comunicar com %s: %w", service, err)
}

func ErrInvalidFormat(field, expected string) error {
	return fmt.Errorf("formato inválido para %s, esperado: %s", field, expected)
}

func ErrInvalidOrderBy(field string) error {
	return fmt.Errorf("'%s' é um campo inválido para ordenação", field)
}

func ErrUnknownWithContext(context string, err error) error {
	return fmt.Errorf("erro desconhecido em %s: %w", context, err)
}
