package handler

import (
	"backend/external/database"
	"backend/external/database/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	fortytwoClientID     = os.Getenv("FORTYTWO_CLIENT_ID")
	fortytwoClientSecret = os.Getenv("FORTYTWO_CLIENT_SECRET")
	fortytwoRedirectURI  = os.Getenv("FORTYTWO_REDIRECT_URI")
)

// AuthWith42 redirects to 42 OAuth authorization page
func AuthWith42(c *gin.Context) {
	if fortytwoClientID == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "42 OAuth not configured",
		})
		return
	}

	// Build authorization URL
	authURL := "https://api.intra.42.fr/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", fortytwoClientID)
	params.Add("redirect_uri", fortytwoRedirectURI)
	params.Add("response_type", "code")

	redirectURL := fmt.Sprintf("%s?%s", authURL, params.Encode())

	fmt.Println("Redirecting to 42 OAuth:", redirectURL)

	c.Redirect(http.StatusFound, redirectURL)
}

// CallbackWith42Complete handles the OAuth callback and exchanges code for token
func CallbackWith42Complete(c *gin.Context) {
	type OAuth42Request struct {
		Code string `json:"code" binding:"required"`
	}

	var req OAuth42Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if fortytwoClientID == "" || fortytwoClientSecret == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "42 OAuth not configured",
		})
		return
	}

	// Exchange code for access token
	tokenURL := "https://api.intra.42.fr/oauth/token"
	tokenData := fmt.Sprintf("grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		fortytwoClientID, fortytwoClientSecret, req.Code, fortytwoRedirectURI)

	fmt.Println("=== Token Exchange Request ===")
	fmt.Println("URL:", tokenURL)
	fmt.Println("Client ID:", fortytwoClientID)
	fmt.Println("Redirect URI:", fortytwoRedirectURI)
	fmt.Println("Code:", req.Code)
	fmt.Println("==============================")

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(tokenData))
	if err != nil {
		fmt.Println("Token exchange error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange code for token",
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("=== Token Exchange Response ===")
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
	fmt.Println("==============================")

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		fmt.Println("Failed to parse token response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse token response",
		})
		return
	}

	if tokenResp.Error != "" {
		fmt.Printf("Token exchange error from 42: %s - %s\n", tokenResp.Error, tokenResp.ErrorDesc)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "42 OAuth error",
			"error_description": tokenResp.ErrorDesc,
		})
		return
	}

	if tokenResp.AccessToken == "" {
		fmt.Println("No access token in response")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No access token received",
			"details": string(body),
		})
		return
	}

	fmt.Println("Access Token received:", tokenResp.AccessToken[:20]+"...")

	// Get user info from 42 API
	userReq, _ := http.NewRequest("GET", "https://api.intra.42.fr/v2/me", nil)
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

	// Debug: Print raw response from 42 API
	fmt.Println("=== 42 API /v2/me Response ===")
	fmt.Println(string(userBody))
	fmt.Println("==============================")

	// 42 API response structure
	var userInfo struct {
		ID            int    `json:"id"`
		Email         string `json:"email"`
		Login         string `json:"login"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		UsualFullName string `json:"usual_full_name"`
		Image         struct {
			Link string `json:"link"`
		} `json:"image"`
	}

	if err := json.Unmarshal(userBody, &userInfo); err != nil {
		fmt.Println("Failed to unmarshal user info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse user info",
			"details": err.Error(),
			"raw":     string(userBody),
		})
		return
	}

	// Debug: Print parsed user info
	fmt.Printf("Parsed User Info: ID=%d, Email=%s, Login=%s, FirstName=%s, LastName=%s\n",
		userInfo.ID, userInfo.Email, userInfo.Login, userInfo.FirstName, userInfo.LastName)

	// Validate required fields
	if userInfo.Email == "" {
		fmt.Println("WARNING: Email is empty! Using login as fallback or checking for alternative email field")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Email not provided by 42 API",
			"user_login": userInfo.Login,
			"user_id":    userInfo.ID,
			"raw_data":   string(userBody),
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
	err = db.Where("email = ? AND strategy = ?", userInfo.Email, "42").First(&user).Error
	if err != nil {
		// Create new user
		firstName := userInfo.FirstName
		if firstName == "" {
			firstName = userInfo.Login
		}

		user = model.User{
			Email:              strings.ToLower(userInfo.Email),
			FirstName:          firstName,
			LastName:           userInfo.LastName,
			Strategy:           "42",
			Language:           "en",
			ProfilePicturePath: userInfo.Image.Link,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
			return
		}
	} else {
		// Update existing user info in case it changed
		user.FirstName = userInfo.FirstName
		if user.FirstName == "" {
			user.FirstName = userInfo.Login
		}
		user.LastName = userInfo.LastName
		user.ProfilePicturePath = userInfo.Image.Link
		db.Save(&user)
	}

	// Generate JWT token with user info
	token, err := GenerateJWTWithUserInfo(user.ID, user.Email, user.FirstName, user.LastName)
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
