package main

import (
	"admin_api/controllers"
	"admin_api/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	db.InitDB()
	router := SetupRouter()
	router.Run(":8083")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("api")
	{
		api.POST("/register", controllers.RegisterStudents)
		api.GET("/commonstudents", controllers.RetrieveCommonStudents)
		api.POST("/suspend", controllers.SuspendStudent)
		api.POST("/retrievefornotifications", controllers.RetrieveStudentsWithNotifications)
	}
	return router
}
