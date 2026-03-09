package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- fake ChatMemberChecker ---

type fakeChatMemberChecker struct {
	isMember bool
}

func (f *fakeChatMemberChecker) IsMember(roomID, userID string) bool {
	return f.isMember
}

// --- fake MessageRepository for Hub ---

// reuse fakeMessageRepo already defined in chat_handler_test.go

// --- Hub unit tests (no WebSocket connection required) ---

func TestNewHub_InitializesFields(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
}

func TestGetOnlineUsers_Empty(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	users := hub.GetOnlineUsers()

	assert.NotNil(t, users)
	assert.Empty(t, users)
}

func TestGetOnlineUsers_WithClients(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	// Add fake clients directly
	c1 := &Client{hub: hub, send: make(chan []byte, 1), userID: "user-1", rooms: make(map[string]bool)}
	c2 := &Client{hub: hub, send: make(chan []byte, 1), userID: "user-2", rooms: make(map[string]bool)}
	c3 := &Client{hub: hub, send: make(chan []byte, 1), userID: "user-1", rooms: make(map[string]bool)} // duplicate

	hub.mu.Lock()
	hub.clients[c1] = true
	hub.clients[c2] = true
	hub.clients[c3] = true
	hub.mu.Unlock()

	users := hub.GetOnlineUsers()

	// Should deduplicate user-1
	assert.Len(t, users, 2)
}

func TestIsUserOnline_False(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	assert.False(t, hub.IsUserOnline("user-1"))
}

func TestIsUserOnline_True(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	c := &Client{hub: hub, send: make(chan []byte, 1), userID: "user-1", rooms: make(map[string]bool)}
	hub.mu.Lock()
	hub.clients[c] = true
	hub.mu.Unlock()

	assert.True(t, hub.IsUserOnline("user-1"))
}

func TestBroadcastToRoom_InvalidJSON(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	// Should return without panicking on bad JSON
	hub.BroadcastToRoom("room-1", []byte("not-json"))
}

func TestBroadcastToRoom_NoClientsInRoom(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	c := &Client{hub: hub, send: make(chan []byte, 1), userID: "user-1", rooms: make(map[string]bool)}
	hub.mu.Lock()
	hub.clients[c] = true
	hub.mu.Unlock()

	msg := `{"type":"message","room_id":"other-room","content":"hello"}`
	// client is not in "other-room", so nothing is sent - should not block or panic
	hub.BroadcastToRoom("other-room", []byte(msg))
}

func TestBroadcastToRoom_ClientInRoom(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})

	c := &Client{hub: hub, send: make(chan []byte, 256), userID: "user-1", rooms: map[string]bool{"room-1": true}}
	hub.mu.Lock()
	hub.clients[c] = true
	hub.mu.Unlock()

	msg := `{"type":"message","room_id":"room-1","content":"hello"}`
	hub.BroadcastToRoom("room-1", []byte(msg))

	// Message should be delivered to the client's send channel
	assert.Len(t, c.send, 1)
}

func TestNewWSHandler_NotNil(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})
	wsh := NewWSHandler(hub)
	assert.NotNil(t, wsh)
}
