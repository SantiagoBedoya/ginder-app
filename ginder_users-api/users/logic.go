package users

import (
	"errors"

	errs "github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound          = errors.New("User not found")
	ErrUserEmailAlreadyExist = errors.New("User email is already in use")
)

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) FindAll() ([]User, error) {
	return s.userRepo.FindAll()
}
func (s *userService) FindOneByID(id string) (*User, error) {
	return s.userRepo.FindOneByID(id)
}
func (s *userService) Create(user *User) (*User, error) {
	exist, err := s.FindOneByEmail(user.Email)
	if err != nil {
		if errs.Cause(err) != ErrUserNotFound {
			return nil, err
		}
	}
	if exist != nil {
		return nil, ErrUserEmailAlreadyExist
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Wrap(err, "Users.Logic.Create")
	}
	user.Password = string(hash)
	return s.userRepo.Create(user)
}
func (s *userService) UpdateOneByID(id string, user *User) error {
	return s.userRepo.UpdateOneByID(id, user)
}
func (s *userService) FindOneByEmail(email string) (*User, error) {
	return s.userRepo.FindOneByEmail(email)
}
func (s *userService) DeleteOneByID(id string) error {
	return s.userRepo.DeleteOneByID(id)
}
