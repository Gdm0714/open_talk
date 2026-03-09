package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/godongmin/open_talk/backend/pkg/response"
)

type chatRoomResponse struct {
	model.ChatRoom
	UnreadCount int64          `json:"unread_count"`
	LastMessage *model.Message `json:"last_message"`
}

type ChatHandler struct {
	chatService service.ChatService
	messageRepo repository.MessageRepository
}

func NewChatHandler(chatService service.ChatService, messageRepo repository.MessageRepository) *ChatHandler {
	return &ChatHandler{chatService: chatService, messageRepo: messageRepo}
}

type createChatRequest struct {
	Type      string   `json:"type" binding:"required"`
	Name      string   `json:"name"`
	MemberIDs []string `json:"member_ids" binding:"required"`
}

type sendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type"`
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID := c.GetString("userID")

	var req createChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	var room *model.ChatRoom
	var err error

	switch model.RoomType(req.Type) {
	case model.RoomTypeDirect:
		if len(req.MemberIDs) != 1 {
			response.BadRequest(c, "direct chat requires exactly one member")
			return
		}
		room, err = h.chatService.CreateDirectChat(userID, req.MemberIDs[0])
	case model.RoomTypeGroup:
		room, err = h.chatService.CreateGroupChat(userID, req.Name, req.MemberIDs)
	default:
		response.BadRequest(c, "invalid chat type: must be 'direct' or 'group'")
		return
	}

	if err != nil {
		response.Error(c, 500, "failed to create chat room")
		return
	}

	response.Created(c, room)
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	userID := c.GetString("userID")

	rooms, err := h.chatService.GetUserChats(userID)
	if err != nil {
		response.Error(c, 500, "failed to get chat rooms")
		return
	}

	result := make([]chatRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		item := chatRoomResponse{ChatRoom: room}

		if count, err := h.messageRepo.GetUnreadCount(room.ID, userID); err == nil {
			item.UnreadCount = count
		}

		if last, err := h.messageRepo.GetLastMessage(room.ID); err == nil {
			item.LastMessage = last
		}

		result = append(result, item)
	}

	response.OK(c, result)
}

func (h *ChatHandler) MarkMessagesRead(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	if err := h.messageRepo.MarkAsRead(roomID, userID); err != nil {
		response.Error(c, 500, "failed to mark messages as read")
		return
	}

	response.Success(c, 200, "messages marked as read")
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.chatService.GetChatMessages(userID, roomID, limit, offset)
	if err != nil {
		if err == service.ErrNotRoomMember {
			response.Error(c, 403, err.Error())
			return
		}
		response.Error(c, 500, "failed to get messages")
		return
	}

	response.OK(c, messages)
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := c.GetString("userID")
	roomID := c.Param("id")

	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if len(req.Content) > 5000 {
		response.BadRequest(c, "message too long (max 5000 characters)")
		return
	}

	msgType := model.MessageTypeText
	if req.MessageType != "" {
		msgType = model.MessageType(req.MessageType)
	}

	msg, err := h.chatService.SendMessage(userID, roomID, req.Content, msgType)
	if err != nil {
		if err == service.ErrNotRoomMember {
			response.Error(c, 403, err.Error())
			return
		}
		response.Error(c, 500, "failed to send message")
		return
	}

	response.Created(c, msg)
}
