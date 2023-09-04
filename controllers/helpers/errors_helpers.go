package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondRecordAlreadyExistError(c *gin.Context, recordType string) {
	data := gin.H{
		"error": recordType + " already exists",
	}
	c.JSON(http.StatusConflict, data)
}

func RespondRecordNotFoundError(c *gin.Context, recordType string) {
	data := gin.H{
		"error": recordType + " not found",
	}
	c.JSON(http.StatusConflict, data)
}

func ResponseInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error. Please contact support."})
}
