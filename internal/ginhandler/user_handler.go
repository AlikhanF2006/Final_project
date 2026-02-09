package ginhandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/AlikhanF2006/Final_project/internal/middleware"
	"github.com/AlikhanF2006/Final_project/internal/postgres/dto"
	"github.com/AlikhanF2006/Final_project/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterDTO
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	u, err := h.svc.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginDTO
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	token, err := h.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Me(c *gin.Context) {
	id := c.GetInt(middleware.UserIDKey)
	u, _ := h.svc.GetProfile(id)
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	id := c.GetInt(middleware.UserIDKey)
	var req dto.UpdateProfileDTO
	c.ShouldBindJSON(&req)
	u, err := h.svc.UpdateProfile(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.GetInt(middleware.UserIDKey)
	var req dto.ChangePasswordDTO
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	if err := h.svc.ChangePassword(id, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) DeleteMe(c *gin.Context) {
	id := c.GetInt(middleware.UserIDKey)
	h.svc.DeleteAccount(id)
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	u, err := h.svc.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) AdminDeleteUser(c *gin.Context) {
	role := c.GetString(middleware.UserRoleKey)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.AdminDeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
