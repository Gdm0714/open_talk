package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// --- Additional auth handler error paths ---

func TestAuthHandler_Register_Returns500OnInternalError(t *testing.T) {
	fake := &fakeAuthService{registerErr: service.ErrInvalidToken} // any non-mapped error
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "ok@example.com", "password": "password123", "nickname": "nick",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAuthHandler_Login_Returns500OnInternalError(t *testing.T) {
	fake := &fakeAuthService{loginErr: service.ErrInvalidToken} // non-mapped error
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/login", gin.H{
		"email": "ok@example.com", "password": "password123",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Additional user handler error paths ---

func TestSearchUsers_ServiceError(t *testing.T) {
	fakeSvc := &fakeUserServiceFull{searchUsersErr: service.ErrRoomNotFound}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	w := getRequest(t, router, "/api/users/search?q=anything")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteAccount_InternalError(t *testing.T) {
	fakeSvc := &fakeUserServiceFull{deleteAccountErr: service.ErrRoomNotFound}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	w := deleteRequest(t, router, "/api/users/account", gin.H{"password": "password123"})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestChangePassword_InternalError(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{
		changePasswordErr: service.ErrRoomNotFound, // non-mapped error -> 500
	})

	w := postJSON(t, router, "/api/users/password", gin.H{
		"old_password": "oldpassword1",
		"new_password": "newpassword1",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Additional friend handler error paths ---

func TestSendFriendRequest_InternalError(t *testing.T) {
	svc := &fakeFriendService{sendRequestErr: service.ErrRoomNotFound}
	router := newFriendRouter(svc)

	w := postJSON(t, router, "/api/friends/request", gin.H{"friend_id": "other-user"})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAcceptFriendRequest_InternalError(t *testing.T) {
	svc := &fakeFriendService{acceptRequestErr: service.ErrRoomNotFound}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/friends/request/fr-1/accept", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRejectFriendRequest_InternalError(t *testing.T) {
	svc := &fakeFriendService{rejectRequestErr: service.ErrRoomNotFound}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/friends/request/fr-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetFriends_ServiceError(t *testing.T) {
	svc := &fakeFriendService{getFriendsErr: service.ErrRoomNotFound}
	router := newFriendRouter(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/friends", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- Additional chat handler error paths ---

func TestGetChats_ServiceError(t *testing.T) {
	chatSvc := &fakeChatService{getUserChatsErr: service.ErrRoomNotFound}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateChat_Direct_ServiceError(t *testing.T) {
	chatSvc := &fakeChatService{createDirectErr: service.ErrRoomNotFound}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats", gin.H{
		"type":       "direct",
		"member_ids": []string{"other-user"},
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateChat_Group_ServiceError(t *testing.T) {
	chatSvc := &fakeChatService{createGroupErr: service.ErrRoomNotFound}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats", gin.H{
		"type":       "group",
		"name":       "Group",
		"member_ids": []string{"user-2"},
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetMessages_InternalError(t *testing.T) {
	chatSvc := &fakeChatService{getChatMessagesErr: service.ErrRoomNotFound}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats/room-1/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSendMessage_InternalError(t *testing.T) {
	chatSvc := &fakeChatService{sendMessageErr: service.ErrRoomNotFound}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats/room-1/messages", gin.H{"content": "hello"})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
