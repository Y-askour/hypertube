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
	githubClientID     = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectURI  = os.Getenv("GITHUB_REDIRECT_URI")
)

// AuthWithGithub handles authentication with GitHub OAuth provider
func AuthWithGithub(c *gin.Context) {
	type OAuthGithubRequest struct {
		Code string `json:"code" binding:"required"`
	}

	var req OAuthGithubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if githubClientID == "" || githubClientSecret == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "GitHub OAuth not configured",
		})
		return
	}

	// Exchange code for access token
	tokenURL := "https://github.com/login/oauth/access_token"
	tokenData := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		githubClientID, githubClientSecret, req.Code, githubRedirectURI)

	tokenReq, _ := http.NewRequest("POST", tokenURL, strings.NewReader(tokenData))
	tokenReq.Header.Set("Accept", "application/json")
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(tokenReq)
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

	// Get user info from GitHub API
	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
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
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(userBody, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse user info",
		})
		return
	}

	// Get email if not in profile (GitHub may require separate API call)
	if userInfo.Email == "" {
		emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		emailReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
		emailResp, err := http.DefaultClient.Do(emailReq)
		if err == nil {
			defer emailResp.Body.Close()
			emailBody, _ := io.ReadAll(emailResp.Body)
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if json.Unmarshal(emailBody, &emails) == nil && len(emails) > 0 {
				for _, e := range emails {
					if e.Primary {
						userInfo.Email = e.Email
						break
					}
				}
				if userInfo.Email == "" {
					userInfo.Email = emails[0].Email
				}
			}
		}
	}

	if userInfo.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to retrieve email from GitHub",
		})
		return
	}

	// Parse name into first and last name
	nameParts := strings.Split(userInfo.Name, " ")
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
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
	err = db.Where("email = ? AND strategy = ?", userInfo.Email, "github").First(&user).Error
	if err != nil {
		// Create new user
		user = model.User{
			Email:     strings.ToLower(userInfo.Email),
			FirstName: firstName,
			LastName:  lastName,
			Strategy:  "github",
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
