package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/godongmin/open_talk/backend/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, token, err := h.authService.Register(req.Email, req.Password, req.Nickname)
	if err != nil {
		switch err {
		case service.ErrInvalidEmail:
			response.BadRequest(c, err.Error())
		case service.ErrInvalidPassword:
			response.BadRequest(c, err.Error())
		case service.ErrInvalidNickname:
			response.BadRequest(c, err.Error())
		case service.ErrEmailAlreadyExists:
			response.Error(c, 409, err.Error())
		default:
			response.Error(c, 500, "internal server error")
		}
		return
	}

	response.Created(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Unauthorized(c)
			return
		}
		response.Error(c, 500, "internal server error")
		return
	}

	response.OK(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	token, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	response.OK(c, gin.H{
		"token": token,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// For JWT-based auth, logout is client-side (delete token)
	// But we return success for API completeness
	response.Success(c, http.StatusOK, "logged out successfully")
}
