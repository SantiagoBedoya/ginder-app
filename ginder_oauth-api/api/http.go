package api

import (
	"log"
	"net/http"

	"github.com/SantiagoBedoya/ginder_oauth-api/dto"
	"github.com/SantiagoBedoya/ginder_oauth-api/oauth"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	SignIn(*gin.Context)
	SignUp(*gin.Context)
}

type handler struct {
	accessTokenService oauth.AccessTokenService
}

func NewAccessTokenHandler(tokenAccessService oauth.AccessTokenService) AccessTokenHandler {
	return &handler{tokenAccessService}
}

func (h *handler) SignIn(c *gin.Context) {
	signIn := &dto.SignIn{}
	if err := c.ShouldBindJSON(signIn); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	at, err := h.accessTokenService.Create(signIn)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, at)
}
func (h *handler) SignUp(c *gin.Context) {
	signUp := &dto.SignUp{}
	if err := c.ShouldBindJSON(signUp); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	at, err := h.accessTokenService.SignUp(signUp)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusCreated, at)
}
