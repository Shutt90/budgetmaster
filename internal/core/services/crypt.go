package services

import (
	"golang.org/x/crypto/bcrypt"
)

type Crypt struct{}

func NewCrypt() Crypt {
	return Crypt{}
}

func (c Crypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (c Crypt) CompareHashAndPassword(password []byte, hash []byte) error {
	return bcrypt.CompareHashAndPassword(password, hash)
}
