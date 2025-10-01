package entity

import (
	"errors"
	"regexp"
	"slices"
	"strings"
)

type Cep string

var (
	ErrZipCodeEmpty         = errors.New("zip code: empty")
	ErrZipCodeInvalidLength = errors.New("zip code: invalid length")
	ErrZipCodeInvalidChars  = errors.New("zip code: invalid characters - only digits allowed")
	ErrZipCodeEqualChars    = errors.New("zip code: eight equal digits not allowed")
)

var zipRegex = regexp.MustCompile(`^[0-9]{8}$`)

var notAllowed = []string{
	"00000000",
	"11111111",
	"22222222",
	"33333333",
	"44444444",
	"55555555",
	"66666666",
	"77777777",
	"88888888",
	"99999999",
}

func NewCep(cep string) (Cep, error) {
	raw := strings.TrimSpace(cep)

	if raw == "" {
		return "", ErrZipCodeEmpty
	}

	if len(raw) != 8 {
		return "", ErrZipCodeInvalidLength
	}

	if !zipRegex.MatchString(raw) {
		return "", ErrZipCodeInvalidChars
	}

	if slices.Contains(notAllowed, raw) {
		return "", ErrZipCodeEqualChars
	}

	return Cep(raw), nil
}

func (c Cep) String() string {
	return string(c)
}

func (c Cep) IsZero() bool {
	return c.String() == ""
}

func (c Cep) Equal(o Cep) bool {
	if c.IsZero() && o.IsZero() {
		return false
	}

	return c.String() == o.String()
}

func (c Cep) MarshalJSON() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Cep) UnmarshalJSON(b []byte) error {
	v, err := NewCep(string(b))
	if err != nil {
		return err
	}

	*c = v
	return nil
}
