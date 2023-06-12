package uuid

import (
	"strings"

	"github.com/google/uuid"
)

func ValidateUUID(src string) (*uuid.UUID, error) {
	u, err := uuid.Parse(strings.ToLower(strings.TrimSpace(src)))
	if err != nil {
		return nil, err
	}
	return &u, nil
}
