package router

import (
	"backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// movies routes
	search := router.Group("movies")
	{
		search.GET("", handler.SearchPopular) // get popular movies with no query if a query is sent get list of movies by query and i can have a sort and filter http querys
		search.GET(":imdbid", nil)            // get movie  details
	}
	return router
}
