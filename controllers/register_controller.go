package controllers

import (
	"admin_api/controllers/helpers"
	"admin_api/db"
	"admin_api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterStudents(c *gin.Context) {
	var requestBody struct {
		Teacher  string   `json:"teacher" binding:"required,email"`
		Students []string `json:"students" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerStudents := []models.Student{}
	newStudents := []models.Student{}

	for _, studentEmail := range requestBody.Students {
		existingStudent, studentExists, err := db.CheckStudentExist(db.DB, studentEmail)
		if err != nil {
			helpers.ResponseInternalServerError(c)
			return
		}
		log.Println(existingStudent, studentExists, err)

		if !studentExists {
			newStudents = append(newStudents, models.Student{Email: studentEmail})
		}
		registerStudents = append(registerStudents, *existingStudent)
	}
	teacherStudentRelations := []models.TeacherStudentRelation{}
	if len(newStudents) > 0 {
		_, err := db.InsertStudents(db.DB, newStudents)
		if err != nil {
			helpers.ResponseInternalServerError(c)
			return
		}
	}

	for _, student := range registerStudents {
		isRegistered, err := db.CheckIfStudentIsRegistered(db.DB, requestBody.Teacher, student.ID)
		if err != nil {
			helpers.ResponseInternalServerError(c)
			return
		}
		if !isRegistered {
			teacherStudentRelations = append(teacherStudentRelations, models.TeacherStudentRelation{
				TeacherEmail: requestBody.Teacher,
				StudentID:    student.ID,
			})
		}
	}

	if err := db.InsertTeacherStudentRelations(db.DB, teacherStudentRelations); err != nil {
		helpers.ResponseInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
