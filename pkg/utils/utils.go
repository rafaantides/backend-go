package utils

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/unicode/norm"
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

func ToStrPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ToFormatDatePointer(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	formatted := t.Format("2006-01-02")
	return &formatted
}

func ToFormatDateTimePointer(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	formated := t.Format("2006-01-02 15:04:05")
	return &formated
}

// TODO: mudar para dar erro depois
// ToUUIDSlice converte um slice de strings para um slice de uuid.UUID, ignorando os inválidos
func ToUUIDSlice(strs []string) []uuid.UUID {
	var result []uuid.UUID
	for _, s := range strs {
		if id, err := uuid.Parse(s); err == nil {
			result = append(result, id)
		}
	}
	return result
}

// ToTimePointer parseia um *string no layout "2006-01-02" para *time.Time
func ToTimePointer(str *string) *time.Time {
	if str == nil || *str == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *str)
	if err != nil {
		return nil
	}
	return &t
}

// removeAccents remove acentos dos caracteres da string
func RemoveAccents(s string) string {
	t := norm.NFD.String(s) // Normaliza os caracteres
	return strings.Map(func(r rune) rune {
		if unicode.IsMark(r) { // Remove caracteres de marcação (acentos)
			return -1
		}
		return r
	}, t)
}

// sanitizeString transforma a string em minúsculas e remove acentos
func SanitizeString(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`) // Remove caracteres especiais, mantendo letras, números e espaço
	return strings.ToLower(re.ReplaceAllString(RemoveAccents(s), ""))
}
