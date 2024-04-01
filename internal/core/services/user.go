package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/Shutt90/budgetmaster/internal/core/ports"
)

type userService struct {
	userRepository ports.UserRepository
}

func NewUserService(ur ports.UserRepository) *userService {
	return &userService{
		userRepository: ur,
	}
}

func (ur *userService) Login(email, password string) error {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = ur.userRepository.GetByLogin(email, passBytes)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userService) ChangePassword(email, password string) error {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = ur.userRepository.ChangePassword(email, string(passBytes))
	if err != nil {
		return err
	}

	return nil
}
