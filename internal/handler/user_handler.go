package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/model"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/service"
)

type Handler struct {
	s *service.Service
}

func NewService(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err})
		return
	}

	user, err := h.s.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) FindAllUsers(c *gin.Context) {
	users, err := h.s.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) FindUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.s.FindUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido"})
		return
	}

	user, err := h.s.UpdateUser(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.s.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuário deletado com sucesso"})
}
