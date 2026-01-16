package handler

import (
	"backend/external/database"
	"backend/external/database/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Register handles user registration with email and password
func Register(c *gin.Context) {
	type RegisterRequest struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Language  string `json:"language"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get database connection
	db := database.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	// Check if user already exists
	var existingUser model.User
	err := db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Set default language if not provided
	language := req.Language
	if language == "" {
		language = "en"
	}

	// Create new user
	newUser := model.User{
		Email:          strings.ToLower(req.Email),
		HashedPassword: hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Strategy:       "email_and_password",
		Language:       language,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(newUser.ID, newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":        newUser.ID,
			"email":     newUser.Email,
			"firstName": newUser.FirstName,
			"lastName":  newUser.LastName,
		},
	})
}
