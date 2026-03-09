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
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- fake AuthService ---

type fakeAuthService struct {
	registerUser  *model.User
	registerToken string
	registerErr   error

	loginUser  *model.User
	loginToken string
	loginErr   error

	refreshToken string
	refreshErr   error

	changePasswordErr error
}

func (f *fakeAuthService) Register(email, password, nickname string) (*model.User, string, error) {
	return f.registerUser, f.registerToken, f.registerErr
}

func (f *fakeAuthService) Login(email, password string) (*model.User, string, error) {
	return f.loginUser, f.loginToken, f.loginErr
}

func (f *fakeAuthService) RefreshToken(tokenString string) (string, error) {
	return f.refreshToken, f.refreshErr
}

func (f *fakeAuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	return f.changePasswordErr
}

// --- fake UserService ---

type fakeUserService struct{}

func (f *fakeUserService) GetProfile(userID string) (*model.User, error) {
	return &model.User{ID: userID}, nil
}

func (f *fakeUserService) UpdateProfile(userID string, nickname, avatarURL, statusMessage *string) (*model.User, error) {
	return &model.User{ID: userID}, nil
}

func (f *fakeUserService) SearchUsers(query string) ([]model.User, error) {
	return nil, nil
}

func (f *fakeUserService) DeleteAccount(userID, password string) error {
	return nil
}

// --- helpers ---

func newTestRouter(h *AuthHandler) *gin.Engine {
	r := gin.New()
	r.POST("/api/auth/register", h.Register)
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/refresh", h.Refresh)
	return r
}

func postJSON(t *testing.T, router *gin.Engine, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var out map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&out))
	return out
}

// --- Register handler tests ---

func TestAuthHandler_Register_Returns201OnSuccess(t *testing.T) {
	fake := &fakeAuthService{
		registerUser:  &model.User{ID: "u1", Email: "ok@example.com", Nickname: "oknick"},
		registerToken: "token-abc",
	}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "ok@example.com", "password": "password123", "nickname": "oknick",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAuthHandler_Register_ResponseBodyContainsUserAndToken(t *testing.T) {
	fake := &fakeAuthService{
		registerUser:  &model.User{ID: "u1", Email: "ok@example.com", Nickname: "oknick"},
		registerToken: "token-abc",
	}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "ok@example.com", "password": "password123", "nickname": "oknick",
	})

	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "token-abc", data["token"])
	assert.NotNil(t, data["user"])
}

func TestAuthHandler_Register_Returns400WhenBodyIsMissingRequiredField(t *testing.T) {
	fake := &fakeAuthService{}
	router := newTestRouter(NewAuthHandler(fake))

	// missing nickname
	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "ok@example.com", "password": "password123",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_Returns400WhenBodyIsEmpty(t *testing.T) {
	fake := &fakeAuthService{}
	router := newTestRouter(NewAuthHandler(fake))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_Returns400OnInvalidEmail(t *testing.T) {
	fake := &fakeAuthService{registerErr: service.ErrInvalidEmail}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "bademail", "password": "password123", "nickname": "nick",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

func TestAuthHandler_Register_Returns400OnInvalidPassword(t *testing.T) {
	fake := &fakeAuthService{registerErr: service.ErrInvalidPassword}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "ok@example.com", "password": "short", "nickname": "nick",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_Returns400OnInvalidNickname(t *testing.T) {
	fake := &fakeAuthService{registerErr: service.ErrInvalidNickname}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "ok@example.com", "password": "password123", "nickname": "x",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_Returns409OnDuplicateEmail(t *testing.T) {
	fake := &fakeAuthService{registerErr: service.ErrEmailAlreadyExists}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/register", gin.H{
		"email": "dup@example.com", "password": "password123", "nickname": "nick",
	})

	assert.Equal(t, http.StatusConflict, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

// --- Login handler tests ---

func TestAuthHandler_Login_Returns200OnSuccess(t *testing.T) {
	fake := &fakeAuthService{
		loginUser:  &model.User{ID: "u2", Email: "login@example.com", Nickname: "loginnick"},
		loginToken: "login-token",
	}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/login", gin.H{
		"email": "login@example.com", "password": "password123",
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Login_ResponseBodyContainsUserAndToken(t *testing.T) {
	fake := &fakeAuthService{
		loginUser:  &model.User{ID: "u2", Email: "login@example.com", Nickname: "loginnick"},
		loginToken: "login-token",
	}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/login", gin.H{
		"email": "login@example.com", "password": "password123",
	})

	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "login-token", data["token"])
}

func TestAuthHandler_Login_Returns401OnWrongCredentials(t *testing.T) {
	fake := &fakeAuthService{loginErr: service.ErrInvalidCredentials}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/login", gin.H{
		"email": "login@example.com", "password": "wrongpass",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

func TestAuthHandler_Login_Returns400WhenBodyIsMissingRequiredField(t *testing.T) {
	fake := &fakeAuthService{}
	router := newTestRouter(NewAuthHandler(fake))

	// missing password
	w := postJSON(t, router, "/api/auth/login", gin.H{"email": "login@example.com"})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Refresh handler tests ---

func TestAuthHandler_Refresh_Returns200OnSuccess(t *testing.T) {
	fake := &fakeAuthService{refreshToken: "new-token"}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/refresh", gin.H{"token": "old-valid-token"})

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "new-token", data["token"])
}

func TestAuthHandler_Refresh_Returns401OnInvalidToken(t *testing.T) {
	fake := &fakeAuthService{refreshErr: service.ErrInvalidToken}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/refresh", gin.H{"token": "bad-token"})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Refresh_Returns400WhenTokenFieldMissing(t *testing.T) {
	fake := &fakeAuthService{}
	router := newTestRouter(NewAuthHandler(fake))

	w := postJSON(t, router, "/api/auth/refresh", gin.H{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- Logout handler tests ---

func newTestRouterWithLogout(h *AuthHandler) *gin.Engine {
	r := gin.New()
	r.POST("/api/auth/logout", h.Logout)
	return r
}

func TestLogout_Success(t *testing.T) {
	fake := &fakeAuthService{}
	router := newTestRouterWithLogout(NewAuthHandler(fake))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

// --- ChangePassword handler tests ---

func newTestRouterWithChangePassword(h *UserHandler) *gin.Engine {
	r := gin.New()
	r.POST("/api/users/password", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.ChangePassword(c)
	})
	return r
}

func TestChangePassword_Handler_Success(t *testing.T) {
	fakeAuth := &fakeAuthService{}
	fakeUser := &fakeUserService{}
	handler := NewUserHandler(fakeUser, fakeAuth)
	router := newTestRouterWithChangePassword(handler)

	w := postJSON(t, router, "/api/users/password", gin.H{
		"old_password": "oldpassword1",
		"new_password": "newpassword1",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestChangePassword_Handler_BadRequest(t *testing.T) {
	fakeAuth := &fakeAuthService{}
	fakeUser := &fakeUserService{}
	handler := NewUserHandler(fakeUser, fakeAuth)
	router := newTestRouterWithChangePassword(handler)

	// missing new_password field
	w := postJSON(t, router, "/api/users/password", gin.H{
		"old_password": "oldpassword1",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}
