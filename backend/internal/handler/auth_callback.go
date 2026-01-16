package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CallbackWith42 handles the 42 OAuth callback - displays code for manual entry
func CallbackWith42(c *gin.Context) {
	code := c.Query("code")
	errorParam := c.Query("error")

	if errorParam != "" {
		c.HTML(http.StatusBadRequest, "callback.html", gin.H{
			"provider": "42",
			"error":    errorParam,
		})
		return
	}

	if code == "" {
		c.HTML(http.StatusBadRequest, "callback.html", gin.H{
			"provider": "42",
			"error":    "No authorization code received",
		})
		return
	}

	// Render the callback page which will automatically exchange the code
	c.HTML(http.StatusOK, "callback.html", gin.H{
		"provider": "42",
		"code":     code,
	})
}

// CallbackWithGithub handles the GitHub OAuth callback
func CallbackWithGithub(c *gin.Context) {
	code := c.Query("code")
	errorParam := c.Query("error")

	if errorParam != "" {
		c.HTML(http.StatusBadRequest, "callback.html", gin.H{
			"provider": "GitHub",
			"error":    errorParam,
		})
		return
	}

	if code == "" {
		c.HTML(http.StatusBadRequest, "callback.html", gin.H{
			"provider": "GitHub",
			"error":    "No authorization code received",
		})
		return
	}

	c.HTML(http.StatusOK, "callback.html", gin.H{
		"provider": "GitHub",
		"code":     code,
	})
}

// CallbackWithGoogle handles the Google OAuth callback
func CallbackWithGoogle(c *gin.Context) {
	code := c.Query("code")
	errorParam := c.Query("error")

	if errorParam != "" {
		c.HTML(http.StatusBadRequest, "callback.html", gin.H{
			"provider": "Google",
			"error":    errorParam,
		})
		return
	}

	if code == "" {
		c.HTML(http.StatusBadRequest, "callback.html", gin.H{
			"provider": "Google",
			"error":    "No authorization code received",
		})
		return
	}

	c.HTML(http.StatusOK, "callback.html", gin.H{
		"provider": "Google",
		"code":     code,
	})
}
