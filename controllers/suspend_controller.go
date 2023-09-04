package controllers

import (
	"admin_api/controllers/helpers"
	"admin_api/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuspendStudent(c *gin.Context) {
	var requestBody struct {
		Student string `json:"student" binding:"required,email"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, studentExists, err := db.CheckStudentExist(db.DB, requestBody.Student)
	if err != nil {
		helpers.ResponseInternalServerError(c)
		return
	}
	if !studentExists {
		helpers.RespondRecordNotFoundError(c, "student")
		return
	}
	if err := db.UpdateStudentSuspensionStatus(db.DB, student); err != nil {
		helpers.ResponseInternalServerError(c)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
