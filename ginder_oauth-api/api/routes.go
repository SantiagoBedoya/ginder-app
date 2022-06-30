package api

import (
	"os"

	"github.com/SantiagoBedoya/ginder_oauth-api/oauth"
	"github.com/SantiagoBedoya/ginder_oauth-api/services/users"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, accessTokenRepo oauth.AccessTokenRepository) {
	userService := users.NewUserService(os.Getenv("USERS_URL"))
	service := oauth.NewAccessTokenService(accessTokenRepo, userService)
	handler := NewAccessTokenHandler(service)
	router.POST("/signup", handler.SignUp)
	router.POST("/signin", handler.SignIn)
}
