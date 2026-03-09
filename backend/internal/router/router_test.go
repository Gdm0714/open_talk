package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.User{},
		&model.ChatRoom{},
		&model.ChatRoomMember{},
		&model.Message{},
		&model.Friend{},
	)
	require.NoError(t, err)
	return db
}

func testConfig() *config.Config {
	return &config.Config{
		Port:      "8080",
		JWTSecret: "test-router-secret",
		DBName:    ":memory:",
	}
}

func generateToken(secret, userID string) string {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte(secret))
	return s
}

func TestSetup_ReturnsNonNilRouter(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()

	r := Setup(db, cfg)

	assert.NotNil(t, r)
}

// --- Public routes should be accessible without auth ---

func TestPublicRoute_Register(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"nickname": "testnick",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPublicRoute_Login(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	// Register first
	regBody, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "password123",
		"nickname": "loginnick",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["data"])
}

func TestPublicRoute_Refresh(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	// Register to get a token
	regBody, _ := json.Marshal(map[string]string{
		"email":    "refresh@example.com",
		"password": "password123",
		"nickname": "refreshnick",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var regResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	data := regResp["data"].(map[string]interface{})
	token := data["token"].(string)

	// Refresh
	refreshBody, _ := json.Marshal(map[string]string{"token": token})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Protected routes should require auth ---

func TestProtectedRoute_GetProfile_NoAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRoute_GetChats_NoAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/chats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRoute_GetFriends_NoAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/friends", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProtectedRoute_WebSocket_NoAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- Protected routes with valid auth ---

func TestProtectedRoute_GetProfile_WithAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	// Register user
	regBody, _ := json.Marshal(map[string]string{
		"email":    "auth@example.com",
		"password": "password123",
		"nickname": "authnick",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var regResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	data := regResp["data"].(map[string]interface{})
	token := data["token"].(string)

	// Access protected route
	req = httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProtectedRoute_GetChats_WithAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	// Register user
	regBody, _ := json.Marshal(map[string]string{
		"email":    "chatauth@example.com",
		"password": "password123",
		"nickname": "chatauthnick",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var regResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	data := regResp["data"].(map[string]interface{})
	token := data["token"].(string)

	req = httptest.NewRequest(http.MethodGet, "/api/chats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProtectedRoute_GetFriends_WithAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	// Register user
	regBody, _ := json.Marshal(map[string]string{
		"email":    "friendauth@example.com",
		"password": "password123",
		"nickname": "friendauthnick",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var regResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	data := regResp["data"].(map[string]interface{})
	token := data["token"].(string)

	req = httptest.NewRequest(http.MethodGet, "/api/friends", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- Method routing ---

func TestRoutes_SearchUsers_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/users/search?q=test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoutes_CreateChat_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	body, _ := json.Marshal(map[string]interface{}{
		"type":       "direct",
		"member_ids": []string{"user-2"},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/chats", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoutes_SendFriendRequest_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	body, _ := json.Marshal(map[string]string{"friend_id": "user-2"})
	req := httptest.NewRequest(http.MethodPost, "/api/friends/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoutes_Logout_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoutes_PasswordChange_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	body, _ := json.Marshal(map[string]string{
		"old_password": "old",
		"new_password": "new",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/users/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoutes_DeleteAccount_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodDelete, "/api/users/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRoutes_MarkRead_RequiresAuth(t *testing.T) {
	db := setupTestDB(t)
	cfg := testConfig()
	r := Setup(db, cfg)

	req := httptest.NewRequest(http.MethodPut, "/api/chats/room-1/read", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
