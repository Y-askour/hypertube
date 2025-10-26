package handler

import (
	"backend/external/omdb"
	"backend/external/yts"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchPopular(c *gin.Context) {
	start := time.Now() // start timer
	var queryParams yts.Params

	if err := c.BindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiResp, err := yts.SearchMovies(queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if apiResp.Data.MovieCount == 0 || len(apiResp.Data.Movies) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no movies found"})
		return
	}

	movies := apiResp.Data.Movies

	// --- Parallel OMDB lookups ---
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // limit to 5 concurrent requests

	for i := range movies {
		wg.Add(1)
		semaphore <- struct{}{} // acquire a slot

		go func(i int) {
			defer wg.Done()
			defer func() { <-semaphore }() // release the slot

			details, err := omdb.SearchMovies(movies[i].ImdbCode)
			if err != nil {
				//fmt.Println("Error fetching OMDB:", err)
				return
			}
			movies[i].Details = details
		}(i)
	}

	wg.Wait()
	// --- End parallel section ---
	fmt.Printf("⏱️ Total OMDB enrichment time: %v\n", time.Since(start))

	apiResp.Data.Movies = movies
	c.JSON(http.StatusOK, apiResp.Data)
}
