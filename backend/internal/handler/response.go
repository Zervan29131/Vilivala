package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// R is a helper to return unified JSON responses
func R(c *gin.Context, httpStatus int, code int, message string, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

// OK sends a successful response
func OK(c *gin.Context, data interface{}) {
	R(c, http.StatusOK, 0, "ok", data)
}

// Fail sends an error response
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	R(c, httpStatus, code, message, nil)
}
