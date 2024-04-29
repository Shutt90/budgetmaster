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

func (c Crypt) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
