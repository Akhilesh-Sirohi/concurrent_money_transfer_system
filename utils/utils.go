package utils

import (
	"github.com/google/uuid"
)

func GenerateUniqueEntityId() string {
	return uuid.New().String()
}
