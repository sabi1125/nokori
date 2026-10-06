//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package util

import "github.com/google/uuid"

type UUIDGenerator interface {
	NewV7() (string, error)
}

type uuidGenerator struct{}

func NewUUIDGenerator() UUIDGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) NewV7() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
