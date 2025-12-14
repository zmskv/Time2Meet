package valueobject

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUID = uuid.UUID

var Nil UUID = uuid.Nil

func ParseUUID(s string) (UUID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Nil, fmt.Errorf("uuid is empty")
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return Nil, fmt.Errorf("invalid uuid: %w", err)
	}
	return u, nil
}
