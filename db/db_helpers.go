package db

import (
	"admin_api/models"

	"gorm.io/gorm"
)

func InsertStudents(db *gorm.DB, students []models.Student) ([]models.Student, error) {
	result := db.Create(&students)
	if result.Error != nil {
		return nil, result.Error
	}
	return students, nil
}

func InsertTeacherStudentRelations(db *gorm.DB, relations []models.TeacherStudentRelation) error {
	result := db.Create(&relations)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CheckIfStudentIsRegistered(db *gorm.DB, teacherEmail string, studentID uint) (bool, error) {
	var count int64
	if err := db.Model(&models.TeacherStudentRelation{}).Where("teacher_email = ? AND student_id = ?", teacherEmail, studentID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckStudentExist(db *gorm.DB, studentEmail string) (*models.Student, bool, error) {
	var student models.Student
	err := db.Where("email=?", studentEmail).First(&student).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &student, true, nil
}

func CheckIfStudentIsSuspended(db *gorm.DB, studentEmail string) (bool, error) {
	var student models.Student
	result := db.Where("email = ?", studentEmail).First(&student)
	if result.Error != nil {
		return false, result.Error
	}
	return student.IsSuspended, nil
}

func UpdateStudentSuspensionStatus(db *gorm.DB, student *models.Student) error {
	return db.Model(&student).Update("IsSuspended", true).Error
}

func GetNonSuspendedStudentsRegisteredUnderTeacher(db *gorm.DB, teacherEmail string) ([]models.Student, error) {
	var students []models.Student
	result := db.Table("students").Select("*").
		Joins("INNER JOIN teacher_student_relations ON students.id = teacher_student_relations").
		Where("teacher_student_relations.teacher_email = ? AND students.id = ?", teacherEmail, false).
		Scan(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	return students, nil
}

func GetCommonStudentsUnderTeachers(db *gorm.DB, teacherEmails []string) ([]models.Student, error) {
	var commonStudents []models.Student
	result := db.Table("students").Select("students.email").
		Joins("INNER JOIN teacher_student_relations ON students.id = teacher_student_relations.student_id").
		Where("teacher_student_relations.teacher_email IN ?", teacherEmails).
		Group("students.email").
		Having("COUNT(DISTINCT teacher_student_relations.teacher_email) = ?", len(teacherEmails)).
		Scan(&commonStudents)
	if result.Error != nil {
		return nil, result.Error
	}

	return commonStudents, nil
}
