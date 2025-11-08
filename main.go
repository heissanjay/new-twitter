package main

import (
	"log"
	"twitter/database"
	"twitter/handlers"
	"twitter/middleware"

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

	routes := router.Group("/")
	{
		routes.POST("/register", handlers.Register)
		routes.POST("/login", handlers.Login)

		routesProtected := routes.Group("/")
		routesProtected.Use(middleware.AuthMiddleware())
		{
			routesProtected.POST("/tweets", handlers.CreateTweet)
			routesProtected.GET("/tweets", handlers.GetTweets)
		}
	}

	router.Run()
}
