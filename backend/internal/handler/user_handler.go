package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/godongmin/open_talk/backend/pkg/response"
)

type UserHandler struct {
	userService service.UserService
	authService service.AuthService
}

func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &UserHandler{userService: userService, authService: authService}
}

type updateProfileRequest struct {
	Nickname      *string `json:"nickname"`
	AvatarURL     *string `json:"avatar_url"`
	StatusMessage *string `json:"status_message"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.NotFound(c)
		return
	}

	response.OK(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, err := h.userService.UpdateProfile(userID, req.Nickname, req.AvatarURL, req.StatusMessage)
	if err != nil {
		response.Error(c, 500, "failed to update profile")
		return
	}

	response.OK(c, user)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "search query is required")
		return
	}

	users, err := h.userService.SearchUsers(query)
	if err != nil {
		response.Error(c, 500, "failed to search users")
		return
	}

	response.OK(c, users)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("userID")

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case service.ErrWrongPassword:
			response.Error(c, http.StatusUnauthorized, err.Error())
		case service.ErrInvalidPassword:
			response.BadRequest(c, err.Error())
		default:
			response.Error(c, 500, "failed to change password")
		}
		return
	}

	response.Success(c, http.StatusOK, "password changed successfully")
}

type deleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("userID")

	var req deleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "password is required")
		return
	}

	if err := h.userService.DeleteAccount(userID, req.Password); err != nil {
		switch err {
		case service.ErrWrongPassword:
			response.Error(c, http.StatusUnauthorized, err.Error())
		default:
			response.Error(c, 500, "failed to delete account")
		}
		return
	}

	response.Success(c, http.StatusOK, "account deleted successfully")
}
