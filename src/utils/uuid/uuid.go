package uuid

import "github.com/google/uuid"

func BuildUuid() string {
	return uuid.New().String()
}
