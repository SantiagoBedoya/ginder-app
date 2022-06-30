package oauth

import (
	"github.com/SantiagoBedoya/ginder_oauth-api/dto"
)

type AccessTokenService interface {
	FindByToken(string) (*AccessToken, error)
	FindByUserID(string) (*AccessToken, error)
	Create(*dto.SignIn) (*AccessToken, error)
	SignUp(*dto.SignUp) (*AccessToken, error)
}
