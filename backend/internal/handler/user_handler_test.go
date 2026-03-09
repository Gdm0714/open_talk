package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// --- configurable fakes for UserHandler tests ---

type fakeUserServiceFull struct {
	getProfileUser *model.User
	getProfileErr  error

	updateProfileUser *model.User
	updateProfileErr  error

	searchUsersResult []model.User
	searchUsersErr    error

	deleteAccountErr error
}

func (f *fakeUserServiceFull) GetProfile(userID string) (*model.User, error) {
	return f.getProfileUser, f.getProfileErr
}

func (f *fakeUserServiceFull) UpdateProfile(userID string, nickname, avatarURL, statusMessage *string) (*model.User, error) {
	return f.updateProfileUser, f.updateProfileErr
}

func (f *fakeUserServiceFull) SearchUsers(query string) ([]model.User, error) {
	return f.searchUsersResult, f.searchUsersErr
}

func (f *fakeUserServiceFull) DeleteAccount(userID, password string) error {
	return f.deleteAccountErr
}

type fakeAuthServiceFull struct {
	changePasswordErr error
}

func (f *fakeAuthServiceFull) Register(email, password, nickname string) (*model.User, string, error) {
	return nil, "", nil
}

func (f *fakeAuthServiceFull) Login(email, password string) (*model.User, string, error) {
	return nil, "", nil
}

func (f *fakeAuthServiceFull) RefreshToken(tokenString string) (string, error) {
	return "", nil
}

func (f *fakeAuthServiceFull) ChangePassword(userID, oldPassword, newPassword string) error {
	return f.changePasswordErr
}

// --- router helpers ---

func newUserRouter(userSvc service.UserService, authSvc service.AuthService) *gin.Engine {
	h := NewUserHandler(userSvc, authSvc)
	r := gin.New()
	r.GET("/api/users/profile", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.GetProfile(c)
	})
	r.PUT("/api/users/profile", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.UpdateProfile(c)
	})
	r.GET("/api/users/search", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.SearchUsers(c)
	})
	r.POST("/api/users/password", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.ChangePassword(c)
	})
	r.DELETE("/api/users/account", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.DeleteAccount(c)
	})
	return r
}

func getRequest(t *testing.T, router *gin.Engine, path string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	router.ServeHTTP(w, req)
	return w
}

func putJSON(t *testing.T, router *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

func deleteRequest(t *testing.T, router *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

// --- GetProfile ---

func TestGetProfile_Success(t *testing.T) {
	fakeSvc := &fakeUserServiceFull{
		getProfileUser: &model.User{ID: "test-user-id", Email: "user@example.com", Nickname: "testnick"},
	}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	w := getRequest(t, router, "/api/users/profile")

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestGetProfile_NotFound(t *testing.T) {
	fakeSvc := &fakeUserServiceFull{
		getProfileErr: service.ErrRoomNotFound,
	}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	w := getRequest(t, router, "/api/users/profile")

	assert.Equal(t, http.StatusNotFound, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

// --- UpdateProfile ---

func TestUpdateProfile_Success(t *testing.T) {
	nick := "newnick"
	fakeSvc := &fakeUserServiceFull{
		updateProfileUser: &model.User{ID: "test-user-id", Nickname: nick},
	}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	w := putJSON(t, router, "/api/users/profile", gin.H{"nickname": nick})

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestUpdateProfile_BadRequest(t *testing.T) {
	fakeSvc := &fakeUserServiceFull{
		updateProfileErr: service.ErrRoomNotFound,
	}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	// service error -> 500
	w := putJSON(t, router, "/api/users/profile", gin.H{"nickname": "anynick"})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- ChangePassword ---

func TestChangePassword_Success(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{})

	w := postJSON(t, router, "/api/users/password", gin.H{
		"old_password": "oldpassword1",
		"new_password": "newpassword1",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestChangePassword_BadRequest(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{})

	// missing new_password
	w := postJSON(t, router, "/api/users/password", gin.H{
		"old_password": "oldpassword1",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

func TestChangePassword_WrongPassword(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{
		changePasswordErr: service.ErrWrongPassword,
	})

	w := postJSON(t, router, "/api/users/password", gin.H{
		"old_password": "wrongoldpass",
		"new_password": "newpassword1",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

// --- DeleteAccount ---

func TestDeleteAccount_Success(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{})

	w := deleteRequest(t, router, "/api/users/account", gin.H{"password": "password123"})

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestDeleteAccount_WrongPassword(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{
		deleteAccountErr: service.ErrWrongPassword,
	}, &fakeAuthServiceFull{})

	w := deleteRequest(t, router, "/api/users/account", gin.H{"password": "wrongpass"})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

func TestDeleteAccount_BadRequest(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{})

	// missing password field -> binding fails
	w := deleteRequest(t, router, "/api/users/account", gin.H{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- SearchUsers ---

func TestSearchUsers_Success(t *testing.T) {
	fakeSvc := &fakeUserServiceFull{
		searchUsersResult: []model.User{
			{ID: "u1", Nickname: "alice"},
			{ID: "u2", Nickname: "bob"},
		},
	}
	router := newUserRouter(fakeSvc, &fakeAuthServiceFull{})

	w := getRequest(t, router, "/api/users/search?q=alice")

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestSearchUsers_MissingQuery(t *testing.T) {
	router := newUserRouter(&fakeUserServiceFull{}, &fakeAuthServiceFull{})

	w := getRequest(t, router, "/api/users/search")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}
