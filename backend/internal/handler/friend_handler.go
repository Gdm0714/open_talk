package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/godongmin/open_talk/backend/pkg/response"
)

type FriendHandler struct {
	friendService service.FriendService
}

func NewFriendHandler(friendService service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

type friendRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
}

func (h *FriendHandler) SendRequest(c *gin.Context) {
	userID := c.GetString("userID")

	var req friendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	friend, err := h.friendService.SendFriendRequest(userID, req.FriendID)
	if err != nil {
		switch err {
		case service.ErrCannotFriendSelf:
			response.BadRequest(c, err.Error())
		case service.ErrFriendRequestExists:
			response.Error(c, 409, err.Error())
		default:
			response.Error(c, 500, "failed to send friend request")
		}
		return
	}

	response.Created(c, friend)
}

func (h *FriendHandler) AcceptRequest(c *gin.Context) {
	userID := c.GetString("userID")
	requestID := c.Param("id")

	friend, err := h.friendService.AcceptFriendRequest(userID, requestID)
	if err != nil {
		if err == service.ErrFriendNotFound {
			response.NotFound(c)
			return
		}
		response.Error(c, 500, "failed to accept friend request")
		return
	}

	response.OK(c, friend)
}

func (h *FriendHandler) RejectRequest(c *gin.Context) {
	userID := c.GetString("userID")
	requestID := c.Param("id")

	if err := h.friendService.RejectFriendRequest(userID, requestID); err != nil {
		if err == service.ErrFriendNotFound {
			response.NotFound(c)
			return
		}
		response.Error(c, 500, "failed to reject friend request")
		return
	}

	response.OK(c, gin.H{"message": "friend request rejected"})
}

func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID := c.GetString("userID")

	friends, err := h.friendService.GetFriends(userID)
	if err != nil {
		response.Error(c, 500, "failed to get friends")
		return
	}

	response.OK(c, friends)
}

// BlockFriend blocks the specified user. The DELETE route previously called
// this handler under the name RemoveFriend but the underlying operation is
// a block, not a deletion — renamed for semantic clarity.
func (h *FriendHandler) BlockFriend(c *gin.Context) {
	userID := c.GetString("userID")
	friendID := c.Param("id")

	if err := h.friendService.BlockUser(userID, friendID); err != nil {
		response.Error(c, 500, "failed to block friend")
		return
	}

	response.OK(c, gin.H{"message": "friend blocked"})
}
