package services

import (
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/Shutt90/budgetmaster/internal/core/domain"
	"github.com/Shutt90/budgetmaster/internal/core/ports"
	"github.com/labstack/gommon/log"
)

type UserService struct {
	userRepository ports.UserRepository
	bcryptIface    ports.Crypt
}

func NewUserService(ur ports.UserRepository, bc ports.Crypt) *UserService {
	return &UserService{
		userRepository: ur,
		bcryptIface:    bc,
	}
}

func (ur *UserService) Login(email, password string) (domain.User, error) {
	u, err := ur.userRepository.GetByEmail(email)
	if err != nil {
		log.Infof("could not find user: %s", email)
		return domain.User{}, err
	}

	if err := ur.bcryptIface.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		log.Errorf("incorrect password for user: %s", email)
		return domain.User{}, errors.New("wrong username/password")
	}

	return u, nil
}

func (ur *UserService) ChangePassword(id, email, password string) error {
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	passBytes, err := ur.bcryptIface.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = ur.userRepository.ChangePassword(idInt, email, string(passBytes))
	if err != nil {
		return err
	}

	return nil
}
