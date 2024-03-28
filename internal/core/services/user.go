package services

import (
	"github.com/Shutt90/budgetmaster/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
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

	err = ur.userRepository.ChangePassword(email, passBytes)
	if err != nil {
		return err
	}

	return nil
}
