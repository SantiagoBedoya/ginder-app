package oauth

type AccessTokenRepository interface {
	FindByToken(string) (*AccessToken, error)
	FindByUserID(string) (*AccessToken, error)
	Create(*AccessToken) (*AccessToken, error)
}
