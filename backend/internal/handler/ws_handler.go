package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type       string `json:"type"`
	RoomID     string `json:"room_id"`
	Content    string `json:"content"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	MessageID  string `json:"message_id,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	userID  string
	rooms   map[string]bool
	roomsMu sync.RWMutex
}

// ChatMemberChecker is the interface Hub uses to verify room membership.
type ChatMemberChecker interface {
	IsMember(roomID, userID string) bool
}

type Hub struct {
	clients     map[*Client]bool
	broadcast   chan []byte
	register    chan *Client
	unregister  chan *Client
	mu          sync.RWMutex
	chatRepo    ChatMemberChecker
	messageRepo repository.MessageRepository
}

func NewHub(chatRepo ChatMemberChecker, messageRepo repository.MessageRepository) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			var wsMsg WSMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				continue
			}

			var stale []*Client
			h.mu.Lock()
			for client := range h.clients {
				client.roomsMu.RLock()
				inRoom := client.rooms[wsMsg.RoomID]
				client.roomsMu.RUnlock()
				if inRoom {
					select {
					case client.send <- message:
					default:
						stale = append(stale, client)
					}
				}
			}
			h.mu.Unlock()

			// Unregister stale clients outside the lock to avoid deadlock.
			for _, client := range stale {
				h.unregister <- client
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		return
	}

	var stale []*Client
	h.mu.Lock()
	for client := range h.clients {
		client.roomsMu.RLock()
		inRoom := client.rooms[roomID]
		client.roomsMu.RUnlock()
		if inRoom {
			select {
			case client.send <- message:
			default:
				stale = append(stale, client)
			}
		}
	}
	h.mu.Unlock()

	// Unregister stale clients outside the lock to avoid deadlock.
	for _, client := range stale {
		h.unregister <- client
	}
}

// GetOnlineUsers returns the list of distinct online user IDs.
func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := make([]string, 0)
	seen := make(map[string]bool)
	for client := range h.clients {
		if !seen[client.userID] {
			seen[client.userID] = true
			users = append(users, client.userID)
		}
	}
	return users
}

// IsUserOnline reports whether at least one connection exists for userID.
func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.userID == userID {
			return true
		}
	}
	return false
}

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	userID := c.GetString("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    h.hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		rooms:  make(map[string]bool),
	}

	h.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		switch wsMsg.Type {
		case "join":
			// Verify room membership before allowing join.
			if c.hub.chatRepo != nil && !c.hub.chatRepo.IsMember(wsMsg.RoomID, c.userID) {
				log.Printf("ws join denied: user %s is not a member of room %s", c.userID, wsMsg.RoomID)
				continue
			}
			c.roomsMu.Lock()
			c.rooms[wsMsg.RoomID] = true
			c.roomsMu.Unlock()

		case "leave":
			c.roomsMu.Lock()
			delete(c.rooms, wsMsg.RoomID)
			c.roomsMu.Unlock()

		case "message":
			// Verify room membership before broadcasting.
			if c.hub.chatRepo != nil && !c.hub.chatRepo.IsMember(wsMsg.RoomID, c.userID) {
				log.Printf("ws message denied: user %s is not a member of room %s", c.userID, wsMsg.RoomID)
				continue
			}

			wsMsg.SenderID = c.userID

			// Persist the message to the database.
			if c.hub.messageRepo != nil {
				msg := &model.Message{
					ChatRoomID:  wsMsg.RoomID,
					SenderID:    c.userID,
					Content:     wsMsg.Content,
					MessageType: model.MessageTypeText,
				}
				if err := c.hub.messageRepo.Create(msg); err != nil {
					log.Printf("ws message persist error: %v", err)
				} else {
					wsMsg.MessageID = msg.ID
					wsMsg.Timestamp = msg.CreatedAt.Format(time.RFC3339)
				}
			}

			data, err := json.Marshal(wsMsg)
			if err != nil {
				continue
			}
			c.hub.broadcast <- data

		case "typing":
			// Broadcast typing indicator to all room members except the sender.
			wsMsg.SenderID = c.userID
			data, err := json.Marshal(wsMsg)
			if err != nil {
				continue
			}

			var stale []*Client
			c.hub.mu.Lock()
			for client := range c.hub.clients {
				client.roomsMu.RLock()
				inRoom := client.rooms[wsMsg.RoomID]
				client.roomsMu.RUnlock()
				if inRoom && client.userID != c.userID {
					select {
					case client.send <- data:
					default:
						stale = append(stale, client)
					}
				}
			}
			c.hub.mu.Unlock()

			for _, sc := range stale {
				c.hub.unregister <- sc
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
