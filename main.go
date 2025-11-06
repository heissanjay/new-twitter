package main

import (
	"log"
	"twitter/database"
	"twitter/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	// initialize DB and run migrations
	err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err = database.RunMigration(); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/tweets", handlers.CreateTweet)
	router.GET("/tweets", handlers.GetTweets)

	router.Run()
}
