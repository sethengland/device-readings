package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var dataMap map[string]map[Reading]struct{}

type Reading struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

type readingBody struct {
	ID       string    `json:"id"`
	Readings []Reading `json:"readings"`
}

func getReading(c *gin.Context) {
	var resp []Reading
	id := c.Param("id")
	if dataMap[id] == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no device with that id"})
		return
	}
	for reading := range dataMap[id] {
		resp = append(resp, reading)
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func postReadings(c *gin.Context) {
	var newReadings readingBody
	if err := c.BindJSON(&newReadings); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "was not able to bind the request body."})
		return
	}
	if dataMap[newReadings.ID] == nil {
		dataMap[newReadings.ID] = make(map[Reading]struct{})
	}
	for _, reading := range newReadings.Readings {
		dataMap[newReadings.ID][reading] = struct{}{}
	}
	var resp []Reading

	for reading := range dataMap[newReadings.ID] {
		resp = append(resp, reading)
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func main() {
	router := gin.Default()
	dataMap = make(map[string]map[Reading]struct{})

	router.POST("/post-readings", postReadings)
	router.GET("/get-reading/:id", getReading)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("exiting with error: " + err.Error())
	}
}
