package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ParseUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format")
	}
	return parsedID, nil
}