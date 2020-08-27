package middleware

import "github.com/google/uuid"

type uuidGenerator struct{}

func newUUIDGenerator() uuidGenerator {
	return uuidGenerator{}
}

func (u uuidGenerator) Run() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
