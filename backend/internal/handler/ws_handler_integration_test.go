package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dialWS connects to a test WebSocket server and returns the connection.
func dialWS(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	return conn
}

// newWSTestServer creates a gin engine + Hub and an httptest.Server.
func newWSTestServer(checker ChatMemberChecker) (*Hub, *httptest.Server) {
	hub := NewHub(checker, &fakeMessageRepo{})
	go hub.Run()

	h := NewWSHandler(hub)
	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		c.Set("userID", "test-user")
		h.HandleWebSocket(c)
	})
	srv := httptest.NewServer(r)
	return hub, srv
}

func TestHandleWebSocket_ConnectAndDisconnect(t *testing.T) {
	hub, srv := newWSTestServer(&fakeChatMemberChecker{isMember: true})
	defer srv.Close()

	conn := dialWS(t, srv)

	// Allow hub to register
	time.Sleep(20 * time.Millisecond)
	assert.True(t, hub.IsUserOnline("test-user"))

	conn.Close()
	time.Sleep(20 * time.Millisecond)
	assert.False(t, hub.IsUserOnline("test-user"))
}

func TestHandleWebSocket_JoinRoom(t *testing.T) {
	hub, srv := newWSTestServer(&fakeChatMemberChecker{isMember: true})
	defer srv.Close()

	conn := dialWS(t, srv)
	defer conn.Close()

	msg := WSMessage{Type: "join", RoomID: "room-1"}
	data, _ := json.Marshal(msg)
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, data))

	time.Sleep(30 * time.Millisecond)

	// User should be subscribed - verify via BroadcastToRoom
	broadcastMsg := WSMessage{Type: "message", RoomID: "room-1", Content: "hello"}
	bData, _ := json.Marshal(broadcastMsg)
	hub.BroadcastToRoom("room-1", bData)

	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, received, err := conn.ReadMessage()
	require.NoError(t, err)

	var got WSMessage
	require.NoError(t, json.Unmarshal(received, &got))
	assert.Equal(t, "hello", got.Content)
}

func TestHandleWebSocket_JoinRoom_Denied(t *testing.T) {
	hub, srv := newWSTestServer(&fakeChatMemberChecker{isMember: false})
	defer srv.Close()

	conn := dialWS(t, srv)
	defer conn.Close()

	msg := WSMessage{Type: "join", RoomID: "room-1"}
	data, _ := json.Marshal(msg)
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, data))

	time.Sleep(30 * time.Millisecond)

	// Not subscribed, so broadcast should not reach client
	broadcastMsg := WSMessage{Type: "message", RoomID: "room-1", Content: "should not arrive"}
	bData, _ := json.Marshal(broadcastMsg)
	hub.BroadcastToRoom("room-1", bData)

	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err := conn.ReadMessage()
	// Should timeout - no message received
	assert.Error(t, err)
}

func TestHandleWebSocket_LeaveRoom(t *testing.T) {
	hub, srv := newWSTestServer(&fakeChatMemberChecker{isMember: true})
	defer srv.Close()

	conn := dialWS(t, srv)
	defer conn.Close()

	// Join then leave
	join := WSMessage{Type: "join", RoomID: "room-1"}
	jData, _ := json.Marshal(join)
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, jData))
	time.Sleep(30 * time.Millisecond)

	leave := WSMessage{Type: "leave", RoomID: "room-1"}
	lData, _ := json.Marshal(leave)
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, lData))
	time.Sleep(30 * time.Millisecond)

	// After leaving, broadcast should not reach client
	broadcastMsg := WSMessage{Type: "message", RoomID: "room-1", Content: "after leave"}
	bData, _ := json.Marshal(broadcastMsg)
	hub.BroadcastToRoom("room-1", bData)

	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err := conn.ReadMessage()
	assert.Error(t, err) // timeout expected
}

func TestHandleWebSocket_SendMessage(t *testing.T) {
	hub, srv := newWSTestServer(&fakeChatMemberChecker{isMember: true})
	defer srv.Close()

	// Connect two clients
	conn1 := dialWS(t, srv)
	defer conn1.Close()

	// Second client with different userID needs a separate route
	hub2 := hub
	r2 := gin.New()
	r2.GET("/ws", func(c *gin.Context) {
		c.Set("userID", "test-user-2")
		NewWSHandler(hub2).HandleWebSocket(c)
	})
	srv2 := httptest.NewServer(r2)
	defer srv2.Close()

	url2 := "ws" + strings.TrimPrefix(srv2.URL, "http") + "/ws"
	conn2, _, err := websocket.DefaultDialer.Dial(url2, nil)
	require.NoError(t, err)
	defer conn2.Close()

	time.Sleep(20 * time.Millisecond)

	// Both join room-1
	joinMsg := WSMessage{Type: "join", RoomID: "room-1"}
	jData, _ := json.Marshal(joinMsg)
	conn1.WriteMessage(websocket.TextMessage, jData)
	conn2.WriteMessage(websocket.TextMessage, jData)
	time.Sleep(30 * time.Millisecond)

	// conn1 sends a message to room-1
	sendMsg := WSMessage{Type: "message", RoomID: "room-1", Content: "hi everyone"}
	sData, _ := json.Marshal(sendMsg)
	conn1.WriteMessage(websocket.TextMessage, sData)

	// conn2 (and conn1 itself) should receive it
	conn2.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, received, err := conn2.ReadMessage()
	require.NoError(t, err)
	var got WSMessage
	json.Unmarshal(received, &got)
	assert.Equal(t, "hi everyone", got.Content)
}

func TestHandleWebSocket_TypingIndicator(t *testing.T) {
	hub, srv := newWSTestServer(&fakeChatMemberChecker{isMember: true})
	defer srv.Close()

	conn1 := dialWS(t, srv)
	defer conn1.Close()

	// Second client on same hub
	r2 := gin.New()
	r2.GET("/ws", func(c *gin.Context) {
		c.Set("userID", "other-user")
		NewWSHandler(hub).HandleWebSocket(c)
	})
	srv2 := httptest.NewServer(r2)
	defer srv2.Close()

	url2 := "ws" + strings.TrimPrefix(srv2.URL, "http") + "/ws"
	conn2, _, err := websocket.DefaultDialer.Dial(url2, nil)
	require.NoError(t, err)
	defer conn2.Close()

	time.Sleep(20 * time.Millisecond)

	// Both join room-1
	joinMsg := WSMessage{Type: "join", RoomID: "room-1"}
	jData, _ := json.Marshal(joinMsg)
	conn1.WriteMessage(websocket.TextMessage, jData)
	conn2.WriteMessage(websocket.TextMessage, jData)
	time.Sleep(30 * time.Millisecond)

	// conn1 sends typing indicator
	typingMsg := WSMessage{Type: "typing", RoomID: "room-1"}
	tData, _ := json.Marshal(typingMsg)
	conn1.WriteMessage(websocket.TextMessage, tData)

	// conn2 should receive the typing indicator
	conn2.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	_, received, err := conn2.ReadMessage()
	require.NoError(t, err)
	var got WSMessage
	json.Unmarshal(received, &got)
	assert.Equal(t, "typing", got.Type)
}

func TestHandleWebSocket_InvalidJSON(t *testing.T) {
	_, srv := newWSTestServer(&fakeChatMemberChecker{isMember: true})
	defer srv.Close()

	conn := dialWS(t, srv)
	defer conn.Close()

	// Send bad JSON - should not crash the server
	conn.WriteMessage(websocket.TextMessage, []byte("not-json"))
	time.Sleep(20 * time.Millisecond)

	// Connection still alive
	ping := WSMessage{Type: "join", RoomID: "room-1"}
	pData, _ := json.Marshal(ping)
	err := conn.WriteMessage(websocket.TextMessage, pData)
	assert.NoError(t, err)
}

func TestHub_Run_RegisterUnregister(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})
	go hub.Run()

	c := &Client{
		hub:    hub,
		send:   make(chan []byte, 1),
		userID: "u1",
		rooms:  make(map[string]bool),
	}

	hub.register <- c
	time.Sleep(20 * time.Millisecond)
	assert.True(t, hub.IsUserOnline("u1"))

	hub.unregister <- c
	time.Sleep(20 * time.Millisecond)
	assert.False(t, hub.IsUserOnline("u1"))
}

func TestHub_Run_Broadcast(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})
	go hub.Run()

	send := make(chan []byte, 4)
	c := &Client{
		hub:    hub,
		send:   send,
		userID: "u1",
		rooms:  map[string]bool{"room-1": true},
	}

	hub.register <- c
	time.Sleep(10 * time.Millisecond)

	msg := WSMessage{Type: "message", RoomID: "room-1", Content: "broadcast test"}
	data, _ := json.Marshal(msg)
	hub.broadcast <- data

	time.Sleep(20 * time.Millisecond)
	assert.Len(t, send, 1)
}

func TestHandleWebSocket_UpgradeFails(t *testing.T) {
	hub := NewHub(&fakeChatMemberChecker{}, &fakeMessageRepo{})
	go hub.Run()
	h := NewWSHandler(hub)

	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		c.Set("userID", "u1")
		h.HandleWebSocket(c)
	})

	// Use a plain HTTP request (no upgrade headers) - upgrade should fail gracefully
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	r.ServeHTTP(w, req)

	// Handler returns without panic; status will be 400 (bad request from upgrader)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
