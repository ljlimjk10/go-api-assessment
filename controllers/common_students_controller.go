package controllers

import (
	"admin_api/controllers/helpers"
	"admin_api/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetrieveCommonStudents(c *gin.Context) {
	teacherEmails := c.QueryArray("teacher")

	if len(teacherEmails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one teacher email must be provided"})
		return
	}

	commonStudents, err := db.GetCommonStudentsUnderTeachers(db.DB, teacherEmails)
	if err != nil {
		helpers.ResponseInternalServerError(c)
		return
	}

	studentEmails := make([]string, len(commonStudents))
	for i, student := range commonStudents {
		studentEmails[i] = student.Email
	}

	c.JSON(http.StatusOK, gin.H{"students": studentEmails})
}
