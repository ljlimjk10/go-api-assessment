package db

import (
	"admin_api/models"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	datasourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", host, user, password, port)
	db, err := gorm.Open(postgres.Open(datasourceName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	result := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "already exists") {
			log.Printf("Database %s already exists, continuing...", dbname)
		} else {
			log.Fatal(result.Error)
		}
	} else {
		log.Printf("Database %s created successfully.", dbname)
	}

	dbSQL, _ := db.DB()
	dbSQL.Close()

	datasourceName = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	DB, err = gorm.Open(postgres.Open(datasourceName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&models.Student{}, &models.TeacherStudentRelation{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migration successful.")
}
