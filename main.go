package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db, err := GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router.GET("/person/:person_id/info", PersonInfoHandler(db))
	router.POST("/person/create", PersonCreateHandler(db))

	log.Println("Server is running on :8080")
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
