package router

import (
	"backend/internal"
	"backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

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

	return router
}
