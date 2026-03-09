package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// --- fake FriendService ---

type fakeFriendService struct {
	sendRequestResult *model.Friend
	sendRequestErr    error

	acceptRequestResult *model.Friend
	acceptRequestErr    error

	rejectRequestErr error

	getFriendsResult []model.Friend
	getFriendsErr    error

	blockUserErr error
}

func (f *fakeFriendService) SendFriendRequest(userID, friendID string) (*model.Friend, error) {
	return f.sendRequestResult, f.sendRequestErr
}

func (f *fakeFriendService) AcceptFriendRequest(userID, requestID string) (*model.Friend, error) {
	return f.acceptRequestResult, f.acceptRequestErr
}

func (f *fakeFriendService) RejectFriendRequest(userID, requestID string) error {
	return f.rejectRequestErr
}

func (f *fakeFriendService) GetFriends(userID string) ([]model.Friend, error) {
	return f.getFriendsResult, f.getFriendsErr
}

func (f *fakeFriendService) BlockUser(userID, friendID string) error {
	return f.blockUserErr
}

// --- router helper ---

func newFriendRouter(svc service.FriendService) *gin.Engine {
	h := NewFriendHandler(svc)
	r := gin.New()
	r.POST("/api/friends/request", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.SendRequest(c)
	})
	r.PUT("/api/friends/request/:id/accept", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.AcceptRequest(c)
	})
	r.DELETE("/api/friends/request/:id", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.RejectRequest(c)
	})
	r.GET("/api/friends", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.GetFriends(c)
	})
	r.DELETE("/api/friends/:id", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.BlockFriend(c)
	})
	return r
}

// --- SendFriendRequest tests ---

func TestSendFriendRequest_Success(t *testing.T) {
	svc := &fakeFriendService{
		sendRequestResult: &model.Friend{ID: "fr-1", UserID: "test-user-id", FriendID: "other-user", Status: model.FriendStatusPending},
	}
	router := newFriendRouter(svc)

	w := postJSON(t, router, "/api/friends/request", gin.H{"friend_id": "other-user"})

	assert.Equal(t, http.StatusCreated, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestSendFriendRequest_Self(t *testing.T) {
	svc := &fakeFriendService{
		sendRequestErr: service.ErrCannotFriendSelf,
	}
	router := newFriendRouter(svc)

	w := postJSON(t, router, "/api/friends/request", gin.H{"friend_id": "test-user-id"})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

func TestSendFriendRequest_AlreadyExists(t *testing.T) {
	svc := &fakeFriendService{
		sendRequestErr: service.ErrFriendRequestExists,
	}
	router := newFriendRouter(svc)

	w := postJSON(t, router, "/api/friends/request", gin.H{"friend_id": "other-user"})

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestSendFriendRequest_BadRequest(t *testing.T) {
	router := newFriendRouter(&fakeFriendService{})

	// missing required friend_id
	w := postJSON(t, router, "/api/friends/request", gin.H{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- AcceptFriendRequest tests ---

func TestAcceptFriendRequest_Success(t *testing.T) {
	svc := &fakeFriendService{
		acceptRequestResult: &model.Friend{ID: "fr-1", Status: model.FriendStatusAccepted},
	}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/friends/request/fr-1/accept", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestAcceptFriendRequest_NotFound(t *testing.T) {
	svc := &fakeFriendService{
		acceptRequestErr: service.ErrFriendNotFound,
	}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/friends/request/nonexistent/accept", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

// --- RejectFriendRequest tests ---

func TestRejectFriendRequest_Success(t *testing.T) {
	router := newFriendRouter(&fakeFriendService{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/friends/request/fr-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestRejectFriendRequest_NotFound(t *testing.T) {
	svc := &fakeFriendService{rejectRequestErr: service.ErrFriendNotFound}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/friends/request/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- GetFriends tests ---

func TestGetFriends_Success(t *testing.T) {
	svc := &fakeFriendService{
		getFriendsResult: []model.Friend{
			{ID: "fr-1", UserID: "test-user-id", FriendID: "user-2", Status: model.FriendStatusAccepted},
			{ID: "fr-2", UserID: "test-user-id", FriendID: "user-3", Status: model.FriendStatusAccepted},
		},
	}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/friends", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
	data := body["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestGetFriends_Empty(t *testing.T) {
	router := newFriendRouter(&fakeFriendService{getFriendsResult: []model.Friend{}})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/friends", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- BlockFriend tests ---

func TestBlockFriend_Success(t *testing.T) {
	router := newFriendRouter(&fakeFriendService{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/friends/user-2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestBlockFriend_Error(t *testing.T) {
	svc := &fakeFriendService{blockUserErr: service.ErrFriendNotFound}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/friends/user-2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
