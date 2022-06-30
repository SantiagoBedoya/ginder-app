package dto

import (
	"errors"
	"strings"
)

type SignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *SignIn) Validate() error {
	if strings.TrimSpace(s.Email) == "" {
		return errors.New("smail should not be empty")
	}
	if strings.TrimSpace(s.Password) == "" {
		return errors.New("password should not be empty")
	}
	return nil
}
