package errs

import (
	"backend-go/pkg/utils"
	"errors"
	"fmt"
	"net/http"
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
	ErrBadRequest         = errors.New(utils.SanitizeString(BadRequest))
	ErrUnauthorized       = errors.New(utils.SanitizeString(Unauthorized))
	ErrForbidden          = errors.New(utils.SanitizeString(Forbidden))
	ErrNotFound           = errors.New(utils.SanitizeString(NotFound))
	ErrConflict           = errors.New(utils.SanitizeString(Conflict))
	ErrUnprocessable      = errors.New(utils.SanitizeString(UnprocessableEntity))
	ErrTooManyRequests    = errors.New(utils.SanitizeString(TooManyRequests))
	ErrInternalServer     = errors.New(utils.SanitizeString(InternalServerError))
	ErrServiceUnavailable = errors.New(utils.SanitizeString(ServiceUnavailable))
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

func InvalidParam(field string, err error) error {
	return fmt.Errorf("%s: %w", field, err)
}

func ResorceNotFound(entity, key string) error {
	return fmt.Errorf("%s não encontrado para %s", entity, key)
}

func ParsingField(field string, err error) error {
	return fmt.Errorf("erro ao processar o campo %s: %w", field, err)
}

func DateParsing(field string) error {
	return fmt.Errorf("erro ao processar o campo %s: use o formato YYYY-MM-DD", field)
}

func UnknownWithContext(context string, err error) error {
	return fmt.Errorf("erro desconhecido em %s: %w", context, err)
}

func FailedToSave(table string, err error) error {
	return fmt.Errorf("failed to save table %s: %w", table, err)
}
