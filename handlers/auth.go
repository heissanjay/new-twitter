package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"twitter/database"
	"twitter/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getJwtSecret())

func getJwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "this-is-a-secret-for-signing-jwt-token"
		log.Print("Unable to JWT secreat from environment variables")
		log.Print("Using default secret")
	}
	return secret
}

func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "unable to generate password hash",
			"details": err.Error(),
		})
		return
	}

	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}
	defer tx.Rollback(ctx)

	var user models.User
	query := `
	INSERT INTO users (username, email, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, username, email, created_at
	`

	now := time.Now()
	err = tx.QueryRow(
		ctx,
		query,
		req.Username,
		req.Email,
		hashedPassword,
		now,
		now,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		if duplicate := strings.Contains(err.Error(), "duplicate key"); duplicate {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "user with this email or user id already exists, please login",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "unable to create user",
			"details": err.Error(),
		})
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate token",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user created",
		"token":   token,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})

		return
	}

	var user models.User
	var passwordHash string

	query := `SELECT id, username, email, password_hash, created_at
			FROM users
			WHERE email = $1
	`

	err := database.DB.QueryRow(context.Background(), query, req.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user does not exists. please register",
			"details": err.Error(),
		})
		return
	}

	// check for password match
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password, please try again",
		})
		return
	}

	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "logged in",
			"token":   token,
		},
	)
}

func generateJWT(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
