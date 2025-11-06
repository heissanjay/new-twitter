package handlers

import (
	"context"
	"net/http"
	"time"
	"twitter/database"
	"twitter/models"

	"github.com/gin-gonic/gin"
)

func CreateTweet(c *gin.Context) {
	var req models.CreateTweetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error":   "invalid request body",
				"details": err.Error(),
			},
		)
		return
	}

	var tweet models.Tweet

	query := `
	INSERT INTO tweets (user_id, content, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, content, created_at, updated_at
	`

	now := time.Now()
	err := database.DB.QueryRow(
		context.Background(),
		query,
		req.UserID,
		req.Content,
		now,
		now,
	).Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt, &tweet.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Unable to create tweet",
			"details": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "tweet created",
			"tweet":   tweet,
		},
	)

}

func GetTweets(c *gin.Context) {
	query := `SELECT id, user_id, content, created_at, updated_at
	FROM tweets
	ORDER by created_at DESC
	LIMIT 100	
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Unable to get tweets",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	tweets := []models.Tweet{}
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt, &tweet.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Unable to fetch tweets",
				"message": err.Error(),
			})

			return
		}
		tweets = append(tweets, tweet)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "tweets fetched",
		"tweets": tweets,
	})
}
