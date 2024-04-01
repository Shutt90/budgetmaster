package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type userService struct {
	userRepository ports.UserRepository
	bcryptIface    ports.Crypt
}

func NewUserService(ur ports.UserRepository, bc ports.Crypt) *userService {
	return &userService{
		userRepository: ur,
		bcryptIface:    bc,
	}
}

func (ur *userService) Login(email, password string) error {
	u, err := ur.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}

	err = ur.bcryptIface.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if err != nil {
		return errors.New("unable to login")
	}
	return nil
}

func (ur *userService) ChangePassword(email, password string) error {
	passBytes, err := ur.bcryptIface.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = ur.userRepository.ChangePassword(email, string(passBytes))
	if err != nil {
		return err
	}

	return nil
}
