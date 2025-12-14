package valueobject

import (
	"fmt"
	"net/mail"
	"strings"
)

type Email string

func ParseEmail(s string) (Email, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("email is empty")
	}
	_, err := mail.ParseAddress(s)
	if err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}
	return Email(strings.ToLower(s)), nil
}

func (e Email) String() string { return string(e) }
