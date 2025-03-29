package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PersonInfoHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		personIDStr := c.Param("person_id")
		personID, err := strconv.Atoi(personIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
			return
		}

		personInfo, err := GetPersonInfo(db, personID)
		if err != nil {
			log.Printf("Error getting person info: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			return
		}

		c.JSON(http.StatusOK, personInfo)
	}
}

func PersonCreateHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var personCreate PersonCreate
		if err := c.ShouldBindJSON(&personCreate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := CreatePerson(db, &personCreate); err != nil {
			log.Printf("Error creating person: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
			return
		}

		c.Status(http.StatusOK)
	}
}
