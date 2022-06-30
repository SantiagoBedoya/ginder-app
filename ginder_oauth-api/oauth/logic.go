package oauth

import (
	"errors"
	"os"
	"time"

	errs "github.com/pkg/errors"

	"github.com/SantiagoBedoya/ginder_oauth-api/dto"
	"github.com/SantiagoBedoya/ginder_oauth-api/services/users"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	accessTokenRepo AccessTokenRepository
	usersService    users.UserService
}

var (
	ErrAccessTokenNotFound = errors.New("AccessToken not found")
)

func NewAccessTokenService(accessTokenRepo AccessTokenRepository, usersService users.UserService) AccessTokenService {
	return &service{accessTokenRepo, usersService}
}

func (s *service) FindByToken(token string) (*AccessToken, error) {
	at, err := s.accessTokenRepo.FindByToken(token)
	if err != nil {
		return nil, err
	}
	if time.Now().Unix() > at.ExpiresAt.Unix() {
		return nil, ErrAccessTokenNotFound
	}
	return at, nil
}

func (s *service) FindByUserID(userID string) (*AccessToken, error) {
	at, err := s.accessTokenRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if time.Now().Unix() > at.ExpiresAt.Unix() {
		return nil, ErrAccessTokenNotFound
	}
	return at, nil
}

func (s *service) Create(signIn *dto.SignIn) (*AccessToken, error) {
	if err := signIn.Validate(); err != nil {
		return nil, err
	}
	user, err := s.usersService.FindByEmail(signIn.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signIn.Password))
	if err != nil {
		return nil, err
	}
	token, err := s.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}
	at, err := s.FindByUserID(user.ID)
	if err != nil {
		if errs.Cause(err) != ErrAccessTokenNotFound {
			return nil, err
		}
	}
	if at != nil {
		return at, nil
	}
	at = &AccessToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
	}
	return s.accessTokenRepo.Create(at)
}

func (s *service) CreateToken(userId string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        userId,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return string(token), err
}

func (s *service) SignUp(data *dto.SignUp) (*AccessToken, error) {
	newUser, err := s.usersService.Create(&users.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	})
	if err != nil {
		return nil, err
	}
	token, err := s.CreateToken(newUser.ID)
	if err != nil {
		return nil, err
	}
	at := &AccessToken{
		Token:     token,
		UserID:    newUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
	}
	return at, nil
}
