package utils

import (
	"regexp"
	"strings"
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
