//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package util

import (
	"backend/internal/domain/apperror"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type passwordHasher struct{}

func NewPasswordHasher() PasswordHasher {
	return &passwordHasher{}
}

func (h *passwordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return "", err
	}
	return string(hash), nil
}

func (h *passwordHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
