package handler

import (
	"backend/external/database"
	"backend/external/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthWithEmailPassword handles authentication with email and password
func AuthWithEmailPassword(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
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

	// Find user by email and strategy
	var user model.User
	err := db.Where("email = ? AND strategy = ?", req.Email, "email_and_password").First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Verify password
	if err := VerifyPassword(user.HashedPassword, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
	})
}
