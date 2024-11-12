package utils

import "github.com/google/uuid"

func GeneratePrefixedUUID(prefix string) string {
	uuid := uuid.New()
	return prefix + "_" + uuid.String()
}
