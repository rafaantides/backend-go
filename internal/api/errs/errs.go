package errs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"details,omitempty"`
}

const (
	BadRequest          = "Requisição inválida"
	Unauthorized        = "Não autorizado"
	Forbidden           = "Acesso proibido"
	NotFound            = "Recurso não encontrado"
	Conflict            = "Conflito nos dados"
	UnprocessableEntity = "Entidade não processável"
	TooManyRequests     = "Muitas requisições, tente novamente mais tarde"
	InternalServerError = "Erro interno no servidor"
	ServiceUnavailable  = "Serviço temporariamente indisponível"
)

var ErrorMessages = map[int]string{
	http.StatusBadRequest:          BadRequest,
	http.StatusUnauthorized:        Unauthorized,
	http.StatusForbidden:           Forbidden,
	http.StatusNotFound:            NotFound,
	http.StatusConflict:            Conflict,
	http.StatusUnprocessableEntity: UnprocessableEntity,
	http.StatusTooManyRequests:     TooManyRequests,
	http.StatusInternalServerError: InternalServerError,
	http.StatusServiceUnavailable:  ServiceUnavailable,
}

var (
	ErrBadRequest         = errors.New(strings.ToLower(BadRequest))
	ErrUnauthorized       = errors.New(strings.ToLower(Unauthorized))
	ErrForbidden          = errors.New(strings.ToLower(Forbidden))
	ErrNotFound           = errors.New(strings.ToLower(NotFound))
	ErrConflict           = errors.New(strings.ToLower(Conflict))
	ErrUnprocessable      = errors.New(strings.ToLower(UnprocessableEntity))
	ErrTooManyRequests    = errors.New(strings.ToLower(TooManyRequests))
	ErrInternalServer     = errors.New(strings.ToLower(InternalServerError))
	ErrServiceUnavailable = errors.New(strings.ToLower(ServiceUnavailable))
)

// APIError representa um erro contendo um status HTTP e uma mensagem
type APIError struct {
	StatusCode int
	Err        error
}

// Implementa a interface error
func (e *APIError) Error() string {
	return e.Err.Error()
}

// NewAPIError cria um novo erro com status code
func NewAPIError(statusCode int, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Err:        err,
	}
}

func AAAInvalidUUID(field string, err error) error {
	return fmt.Errorf("%s: %w", field, err)
}

func AAANotFound(entity, key string) error {
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

