package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Success: false,
		Error:   msg,
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error:   "unauthorized",
	})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "not found",
	})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   msg,
	})
}

func Success(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Success: true,
		Message: msg,
	})
}
