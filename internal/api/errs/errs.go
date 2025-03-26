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
	ErrInvalidField     = errors.New("campo inválido")
	ErrInvalidPage      = errors.New("página inválida")
	ErrInvalidPageSize  = errors.New("tamanho da página inválido")
)

func NotFound(entity, key string) error {
	return fmt.Errorf("%s não encontrado para %s", entity, key)
}

func ParsingField(field string, err error) error {
	return fmt.Errorf("erro ao processar o campo %s: %w", field, err)
}

func DateParsing(field string) error {
	return fmt.Errorf("erro ao processar o campo %s: use o formato YYYY-MM-DD", field)
}

func InvalidOrderBy(field string) error {
	return fmt.Errorf("'%s' é um valor inválido para ordenação", field)
}

func UnknownWithContext(context string, err error) error {
	return fmt.Errorf("erro desconhecido em %s: %w", context, err)
}

func InvalidUUID(field string, err string) string {
	return fmt.Sprintf("%s %s", field, err)
}

// func ErrFetchingData(entity, key string, err error) error {
// 	return fmt.Errorf("erro ao buscar %s para %s: %w", entity, key, err)
// }

// func ErrConflict(entity, key string) error {
// 	return fmt.Errorf("conflito ao processar %s com %s", entity, key)
// }

// func ErrExternalDependency(service string, err error) error {
// 	return fmt.Errorf("erro ao comunicar com %s: %w", service, err)
// }

// func ErrInvalidFormat(field, expected string) error {
// 	return fmt.Errorf("formato inválido para %s, esperado: %s", field, expected)
// }
