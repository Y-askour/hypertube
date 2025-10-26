package omdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const omdbBaseUrl = "http://www.omdbapi.com?apikey="
const apiKey = "76567ea7"

// OMDBRating represents a single rating source/value pair
type OMDBRating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// OMDBResult represents a single movie result from the OMDb API
type OMDBResult struct {
	Title      string       `json:"Title"`
	Year       string       `json:"Year"`
	Rated      string       `json:"Rated"`
	Released   string       `json:"Released"`
	Runtime    string       `json:"Runtime"`
	Genre      string       `json:"Genre"`
	Director   string       `json:"Director"`
	Writer     string       `json:"Writer"`
	Actors     string       `json:"Actors"`
	Plot       string       `json:"Plot"`
	Language   string       `json:"Language"`
	Country    string       `json:"Country"`
	Awards     string       `json:"Awards"`
	Poster     string       `json:"Poster"`
	Ratings    []OMDBRating `json:"Ratings"`
	Metascore  string       `json:"Metascore"`
	ImdbRating string       `json:"imdbRating"`
	ImdbVotes  string       `json:"imdbVotes"`
	ImdbID     string       `json:"imdbID"`
	Type       string       `json:"Type"`
	DVD        string       `json:"DVD"`
	BoxOffice  string       `json:"BoxOffice"`
	Production string       `json:"Production"`
	Website    string       `json:"Website"`
	Response   string       `json:"Response"`
}

func SearchMovies(imdbID string) (*OMDBResult, error) {
	path := fmt.Sprintf("%s%s&i=%s", omdbBaseUrl, apiKey, imdbID)
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp OMDBResult
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Response == "False" {
		return nil, fmt.Errorf("OMDB: no result for %s", imdbID)
	}

	return &apiResp, nil
}
