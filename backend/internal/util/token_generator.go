//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"backend/internal/domain/apperror"
)

type TokenGenerator interface {
	GenerateToken() (token string, tokenHash string, err error)
}

type tokenGenerator struct{}

func NewTokenGenerator() TokenGenerator {
	return &tokenGenerator{}
}

func (tokenGenerator *tokenGenerator) GenerateToken() (token string, tokenHash string, err error) {
	bytes := make([]byte, 32)

	if _, err = rand.Read(bytes); err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}

	token = base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(token))
	tokenHash = hex.EncodeToString(hash[:])

	return
}
