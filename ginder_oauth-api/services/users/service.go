package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type UserService interface {
	Create(*User) (*User, error)
	FindByEmail(string) (*User, error)
}

type service struct {
	baseURL string
}

func NewUserService(baseURL string) UserService {
	return &service{baseURL}
}
func (s *service) Create(user *User) (*User, error) {
	fmt.Println("making request to: ", s.baseURL)
	jsonData, _ := json.Marshal(user)
	request, err := http.Post(s.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "Services.User.Create")
	}
	newUser := &User{}
	if err := json.NewDecoder(request.Body).Decode(newUser); err != nil {
		return nil, errors.Wrap(err, "Services.User.Create")
	}
	fmt.Println(newUser)
	return newUser, nil
}
func (s *service) FindByEmail(email string) (*User, error) {
	url := fmt.Sprintf("%s/by-email", s.baseURL)
	jsonData, _ := json.Marshal(map[string]string{"email": email})
	request, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "Services.User.FindByEmail")
	}
	user := &User{}
	if err := json.NewDecoder(request.Body).Decode(user); err != nil {
		return nil, errors.Wrap(err, "Services.User.FindByEmail")
	}
	return user, nil
}
