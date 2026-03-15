package handler

import "github.com/gin-gonic/gin"

func NewConnection(r *gin.Engine, handler *Handler) {
	users := r.Group("/api")

	users.POST("/user", handler.CreateUser)
	users.GET("/user", handler.FindAllUsers)
	users.GET("/user/:id", handler.FindUserById)
	users.PUT("/user/:id", handler.Update)
	users.DELETE("user/:id", handler.Delete)
}
