package controllers

import (
	"admin_api/controllers/helpers"
	"admin_api/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetrieveStudentsWithNotifications(c *gin.Context) {
	var requestBody struct {
		Teacher      string `json:"teacher" binding:"required,email"`
		Notification string `json:"notification" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	students, err := db.GetNonSuspendedStudentsRegisteredUnderTeacher(db.DB, requestBody.Teacher)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	registeredStudentEmails := make([]string, len(students))
	for i, student := range students {
		registeredStudentEmails[i] = student.Email
	}

	mentionedStudentEmails := helpers.ExtractMentionedStudentEmails(requestBody.Notification)
	log.Println(mentionedStudentEmails)
	nonSuspendedMentionedStudentEmails := make([]string, 0)

	if len(mentionedStudentEmails) > 0 {
		for _, studentEmail := range mentionedStudentEmails {
			isSuspended, err := db.CheckIfStudentIsSuspended(db.DB, studentEmail)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if !isSuspended {
				nonSuspendedMentionedStudentEmails = append(nonSuspendedMentionedStudentEmails, studentEmail)
			}
		}
	}

	finalStudentEmails := helpers.RemoveDuplicates(registeredStudentEmails, nonSuspendedMentionedStudentEmails)
	c.JSON(http.StatusOK, gin.H{"recipients": finalStudentEmails})
}
