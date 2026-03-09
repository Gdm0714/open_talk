package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func testCfg() *config.Config {
	return &config.Config{JWTSecret: "test-secret-key"}
}

func validToken(secret, userID string) string {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte(secret))
	return s
}

func expiredToken(secret, userID string) string {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := tok.SignedString([]byte(secret))
	return s
}

func setupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(AuthMiddleware(cfg))
	r.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return r
}

func TestAuthMiddleware_ValidBearerToken(t *testing.T) {
	cfg := testCfg()
	r := setupRouter(cfg)
	tok := validToken(cfg.JWTSecret, "user-123")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user-123")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	r := setupRouter(testCfg())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	r := setupRouter(testCfg())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NonBearerScheme(t *testing.T) {
	cfg := testCfg()
	r := setupRouter(cfg)
	tok := validToken(cfg.JWTSecret, "user-123")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	cfg := testCfg()
	r := setupRouter(cfg)
	tok := expiredToken(cfg.JWTSecret, "user-123")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	cfg := testCfg()
	r := setupRouter(cfg)
	tok := validToken("wrong-secret", "user-123")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MalformedToken(t *testing.T) {
	r := setupRouter(testCfg())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt.token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_QueryParamFallback(t *testing.T) {
	cfg := testCfg()
	r := setupRouter(cfg)
	tok := validToken(cfg.JWTSecret, "ws-user-42")

	req := httptest.NewRequest(http.MethodGet, "/protected?token="+tok, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ws-user-42")
}

func TestAuthMiddleware_HeaderTakesPrecedenceOverQuery(t *testing.T) {
	cfg := testCfg()
	r := setupRouter(cfg)
	headerTok := validToken(cfg.JWTSecret, "header-user")
	queryTok := validToken(cfg.JWTSecret, "query-user")

	req := httptest.NewRequest(http.MethodGet, "/protected?token="+queryTok, nil)
	req.Header.Set("Authorization", "Bearer "+headerTok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "header-user")
}

func TestAuthMiddleware_SetsUserIDInContext(t *testing.T) {
	cfg := testCfg()
	tok := validToken(cfg.JWTSecret, "ctx-user-99")

	var capturedUserID string
	r := gin.New()
	r.Use(AuthMiddleware(cfg))
	r.GET("/check", func(c *gin.Context) {
		capturedUserID = c.GetString("userID")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ctx-user-99", capturedUserID)
}

func TestAuthMiddleware_EmptyBearerToken(t *testing.T) {
	r := setupRouter(testCfg())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_NonHMACSigningMethod(t *testing.T) {
	// Create a token with none algorithm (should be rejected)
	claims := jwt.RegisteredClaims{
		Subject:   "user-123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, _ := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)

	r := setupRouter(testCfg())
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
