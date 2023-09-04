package tests

import (
	"admin_api/db"
	"admin_api/models"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	TDB *gorm.DB
)

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	TDB = setupTestDB()

	os.Exit(m.Run())
}

func setupTestDB() *gorm.DB {
	testHost := os.Getenv("TEST_DB_HOST")
	testUser := os.Getenv("TEST_DB_USER")
	testPassword := os.Getenv("TEST_DB_PASSWORD")
	testDBName := os.Getenv("TEST_DB_NAME")
	testPort := os.Getenv("TEST_DB_PORT")

	datasourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		testHost,
		testUser,
		testPassword,
		testPort,
	)

	db, err := gorm.Open(postgres.Open(datasourceName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	result := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDBName))
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "already exists") {
			log.Printf("Database %s already exists, continuing...", testDBName)
		} else {
			log.Fatal(result.Error)
		}
	} else {
		log.Printf("Database %s created successfully.", testDBName)
	}

	dbSQL, _ := db.DB()
	dbSQL.Close()

	datasourceName = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		testHost,
		testUser,
		testPassword,
		testDBName,
		testPort,
	)
	db, err = gorm.Open(postgres.Open(datasourceName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func refreshTables() error {
	err := TDB.Migrator().DropTable("students")
	if err != nil {
		return err
	}
	err = TDB.AutoMigrate(&models.Student{}, &models.TeacherStudentRelation{})
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table(s)")
	return nil
}

func insertStudent() ([]models.Student, error) {
	student := []models.Student{
		{
			Email:       "studentjon@gmail.com",
			IsSuspended: false,
		},
	}
	err := TDB.Model(&models.Student{}).Create(&student).Error
	if err != nil {
		return nil, err
	}
	return student, nil
}

func insertStudents() ([]models.Student, error) {
	students := []models.Student{
		{
			Email:       "studentjon@gmail.com",
			IsSuspended: false,
		},
		{
			Email:       "studenthon@gmail.com",
			IsSuspended: false,
		},
	}
	for i := range students {
		err := TDB.Model(&models.Student{}).Create(&students[i]).Error
		if err != nil {
			return nil, err
		}
	}
	return students, nil
}

func insertSuspendedStudent() ([]models.Student, error) {
	student := []models.Student{
		{
			Email:       "studentjon@gmail.com",
			IsSuspended: true,
		},
	}
	err := TDB.Model(&models.Student{}).Create(&student).Error
	if err != nil {
		return nil, err
	}
	return student, nil
}

func TestInsertStudent(t *testing.T) {
	err := refreshTables()
	if err != nil {
		t.Fatal(err)
	}

	_, _ = insertStudent()
	studentCount := 0
	_, student1Exist, _ := db.CheckStudentExist(TDB, "studentjon@gmail.com")
	if student1Exist {
		studentCount += 1
	}
	assert.Equal(t, studentCount, 1)
}

func TestInsertStudents(t *testing.T) {
	err := refreshTables()
	if err != nil {
		t.Fatal(err)
	}

	_, err = insertStudents()
	if err != nil {
		t.Fatal(err)
	}
	studentCount := 0
	_, student1Exist, err := db.CheckStudentExist(TDB, "studentjon@gmail.com")
	if err != nil {
		t.Fatal(err)
	}
	if student1Exist {
		studentCount += 1
	}
	_, student2Exist, err := db.CheckStudentExist(TDB, "studenthon@gmail.com")
	if err != nil {
		t.Fatal(err)
	}
	if student2Exist {
		studentCount += 1
	}
	assert.Equal(t, studentCount, 2)
}

func TestSuspendStudent(t *testing.T) {
	err := refreshTables()
	if err != nil {
		t.Fatal(err)
	}

	_, err = insertSuspendedStudent()
	if err != nil {
		t.Fatal(err)
	}
	studentCount := 0
	student1Exist, _ := db.CheckIfStudentIsSuspended(TDB, "studentjon@gmail.com")
	if student1Exist {
		studentCount += 1
	}
	assert.Equal(t, studentCount, 1)

}
