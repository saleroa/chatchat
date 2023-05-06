package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseSuccess(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": message,
	})
}
func ResponseFail(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  500,
		"message": message,
	})
}
