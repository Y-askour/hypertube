package handler

import (
	"backend/external/database"
	"backend/external/database/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	googleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI  = os.Getenv("GOOGLE_REDIRECT_URI")
)

// AuthWithGoogle handles authentication with Google OAuth provider
func AuthWithGoogle(c *gin.Context) {
	type OAuthGoogleRequest struct {
		Code string `json:"code" binding:"required"`
	}

	var req OAuthGoogleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if googleClientID == "" || googleClientSecret == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Google OAuth not configured",
		})
		return
	}

	// Exchange code for access token
	tokenURL := "https://oauth2.googleapis.com/token"
	tokenData := fmt.Sprintf("grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		googleClientID, googleClientSecret, req.Code, googleRedirectURI)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(tokenData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange code for token",
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse token response",
		})
		return
	}

	// Get user info from Google API
	userReq, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user info",
		})
		return
	}
	defer userResp.Body.Close()

	userBody, _ := io.ReadAll(userResp.Body)
	var userInfo struct {
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}
	if err := json.Unmarshal(userBody, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse user info",
		})
		return
	}

	// Find or create user in database
	db := database.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var user model.User
	err = db.Where("email = ? AND strategy = ?", userInfo.Email, "google").First(&user).Error
	if err != nil {
		// Create new user
		user = model.User{
			Email:     strings.ToLower(userInfo.Email),
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Strategy:  "google",
			Language:  "en",
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}
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
