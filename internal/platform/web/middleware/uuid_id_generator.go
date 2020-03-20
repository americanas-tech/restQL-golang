package middleware

import "github.com/google/uuid"

type UuidIdGenerator struct{}

func NewUuidIdGenerator() UuidIdGenerator {
	return UuidIdGenerator{}
}

func (u UuidIdGenerator) Run() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
