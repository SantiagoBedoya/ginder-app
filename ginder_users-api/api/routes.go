package api

import (
	"github.com/SantiagoBedoya/ginder_users-api/users"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userRepo users.UserRepository) {
	service := users.NewUserService(userRepo)
	handler := NewUserHandler(service)

	router.GET("/", handler.GetAll)
	router.POST("/", handler.Create)
	router.POST("/by-email", handler.GetOneByEmail)
	router.GET("/:id", handler.GetOneByID)
	router.PUT("/:id", handler.UpdateOneByID)
	router.DELETE("/:id", handler.DeleteOneByID)
}
