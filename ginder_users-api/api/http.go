package api

import (
	"log"
	"net/http"

	"github.com/SantiagoBedoya/ginder_users-api/users"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UserHandler interface {
	GetAll(*gin.Context)
	GetOneByID(*gin.Context)
	GetOneByEmail(*gin.Context)
	Create(*gin.Context)
	UpdateOneByID(*gin.Context)
	DeleteOneByID(*gin.Context)
}

type handler struct {
	userService users.UserService
}

func NewUserHandler(userService users.UserService) UserHandler {
	return &handler{userService}
}

func (h *handler) GetOneByEmail(c *gin.Context) {
	user := &users.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	currentUser, err := h.userService.FindOneByEmail(user.Email)
	if err != nil {
		if errors.Cause(err) == users.ErrUserNotFound {
			c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, currentUser)
}

func (h *handler) GetAll(c *gin.Context) {
	users, err := h.userService.FindAll()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, users)
}
func (h *handler) GetOneByID(c *gin.Context) {
	user, err := h.userService.FindOneByID(c.Param("id"))
	if err != nil {
		if errors.Cause(err) == users.ErrUserNotFound {
			c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusOK, user)
}
func (h *handler) Create(c *gin.Context) {
	user := &users.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	createdUser, err := h.userService.Create(user)
	if err != nil {
		if errors.Cause(err) == users.ErrUserEmailAlreadyExist {
			c.JSON(http.StatusBadRequest, users.ErrUserEmailAlreadyExist.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}
func (h *handler) UpdateOneByID(c *gin.Context) {
	user := &users.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	err := h.userService.UpdateOneByID(c.Param("id"), user)
	if err != nil {
		if errors.Cause(err) == users.ErrUserNotFound {
			c.JSON(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *handler) DeleteOneByID(c *gin.Context) {
	err := h.userService.DeleteOneByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.Status(http.StatusNoContent)
}
