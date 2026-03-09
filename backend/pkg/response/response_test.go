package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, nil)
	c.Request = req
	return c, w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	t.Helper()
	var resp Response
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	return resp
}

func TestOK_ReturnsStatusOK(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	OK(c, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOK_ReturnsSuccessTrue(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	OK(c, nil)
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
}

func TestOK_IncludesData(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	OK(c, gin.H{"key": "value"})
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "value", data["key"])
}

func TestCreated_ReturnsStatusCreated(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	Created(c, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreated_ReturnsSuccessTrue(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	Created(c, gin.H{"id": "123"})
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
}

func TestError_ReturnsProvidedStatusCode(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	Error(c, http.StatusInternalServerError, "something went wrong")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestError_ReturnsSuccessFalse(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	Error(c, http.StatusInternalServerError, "something went wrong")
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
}

func TestError_IncludesErrorMessage(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	Error(c, http.StatusConflict, "email already exists")
	resp := decodeResponse(t, w)
	assert.Equal(t, "email already exists", resp.Error)
}

func TestUnauthorized_ReturnsStatusUnauthorized(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	Unauthorized(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUnauthorized_ReturnsSuccessFalse(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	Unauthorized(c)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, "unauthorized", resp.Error)
}

func TestNotFound_ReturnsStatusNotFound(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	NotFound(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNotFound_ReturnsSuccessFalse(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/")
	NotFound(c)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, "not found", resp.Error)
}

func TestBadRequest_ReturnsStatusBadRequest(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	BadRequest(c, "invalid input")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBadRequest_ReturnsSuccessFalse(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	BadRequest(c, "invalid input")
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
}

func TestBadRequest_IncludesErrorMessage(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	BadRequest(c, "missing required field")
	resp := decodeResponse(t, w)
	assert.Equal(t, "missing required field", resp.Error)
}

func TestSuccess_ReturnsCorrectStatusCodeAndMessage(t *testing.T) {
	c, w := newTestContext(http.MethodPost, "/")
	Success(c, http.StatusOK, "operation successful")
	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, "operation successful", resp.Message)
	assert.Empty(t, resp.Error)
}
