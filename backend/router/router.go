package router

import (
	"backend/internal"
	"backend/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// subtitle package
	"backend/subtitle"
)

func Router() *gin.Engine {
	router := gin.Default()

	// Enable CORS for frontend
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// movies routes
	movies := router.Group("movie")
	{

		movies.GET("popular", handler.SearchPopular) // get popular movies with no query if a query is sent get list of movies by query and i can have a sort and filter http querys
		movies.GET("search", nil)                    // get popular movies with no query if a query is sent get list of movies by query and i can have a sort and filter http querys
		movies.GET(":imdbid", nil)                   // get movie  details
	}

	comment := router.Group("comment")
	{
		comment.GET("", nil)
	}

	stream := router.Group("stream")
	{
		stream.GET("", internal.Stream)
	}

	// auth routes
	auth := router.Group("auth")
	{
		// OAuth providers - redirect to authorization pages
		auth.GET("42", handler.AuthWith42)                             // authenticate with 42
		auth.GET("42/callback", handler.CallbackWith42)                // 42 OAuth callback
		auth.POST("42/callback", handler.CallbackWith42Complete)       // 42 code exchange endpoint
		auth.POST("github", handler.AuthWithGithub)                    // authenticate with GitHub
		auth.GET("github/callback", handler.CallbackWithGithub)        // GitHub OAuth callback
		auth.POST("google", handler.AuthWithGoogle)                    // authenticate with Google
		auth.GET("google/callback", handler.CallbackWithGoogle)        // Google OAuth callback
		auth.POST("email_and_password", handler.AuthWithEmailPassword) // authenticate with email and password
		auth.POST("register", handler.Register)                        // register new user
	}

	// subtitle routes
	subtitleGroup := router.Group("subtitle")
	{	
		// search for subtitles by query (movie name, imdb id, etc.)
		subtitleGroup.GET("/search", subtitle.SearchSubtitlesHandler)

		// get download link for a subtitle by file id and optional parameters (format, filename, force download)
		// subtitleGroup.POST("/download", subtitle.GetDownloadLinkHandler)

		subtitleGroup.POST("/download", subtitle.DownloadSubtitleHandler)

		// download subtitle file by file id (this will be a proxy endpoint that fetches the file from OpenSubtitles and serves it to the client)
		// subtitleGroup.GET("/file", subtitle.DownloadSubtitleFileHandler)
	}

	return router
}
