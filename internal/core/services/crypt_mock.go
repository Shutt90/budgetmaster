package services

import "errors"

var ErrNoMatch = errors.New("error not matching")

type mockCrypt struct{}

func NewMockCrypt() mockCrypt {
	return mockCrypt{}
}

func (c mockCrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return []byte("password"), nil
}

func (c mockCrypt) CompareHashAndPassword(password []byte, hash []byte) error {
	hash = []byte("password")
	if string(password) == string("password") {
		return nil
	}

	return ErrNoMatch
}
