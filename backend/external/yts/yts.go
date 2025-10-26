package yts

import (
	"backend/external/omdb"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var path = "https://yts.mx/api/v2/list_movies.json"

type Params struct {
	Limit      uint   `json:"limit" form:"limit"`
	Page       uint   `json:"page" form:"page"`
	Quality    string `json:"quality" form:"quality"` //  (480p, 720p, 1080p, 1080p.x265, 2160p, 3D)
	Query_term string `json:"query_term" form:"query_term"`
	Genre      string `json:"genre" form:"genre"`
	Sort_by    string `json:"sort_by" form:"sort_by"` //  (title, year, rating, peers, seeds, download_count, like_count, date_added)
}

func (p Params) ToQueryParams() string {
	values := url.Values{}
	if p.Limit != 0 {
		values.Add("limit", strconv.Itoa(int(p.Limit)))
	}
	if p.Page != 0 {
		values.Add("page", strconv.Itoa(int(p.Page)))
	}
	if p.Quality != "" {
		values.Add("quality", p.Quality)
	}
	if p.Query_term != "" {
		values.Add("query_term", p.Query_term)
	}
	if p.Genre != "" {
		values.Add("genre", p.Genre)
	}
	if p.Sort_by != "" {
		values.Add("sort_by", p.Sort_by)
	}
	return values.Encode()
}

type APIResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	MovieCount int     `json:"movie_count"`
	Movies     []Movie `json:"movies"`
}
type Movie struct {
	ID        int              `json:"id"`
	Title     string           `json:"title"`
	TitleLong string           `json:"title_long"`
	ImdbCode  string           `json:"imdb_code"`
	Torrents  []Torrent        `json:"torrents"`
	Details   *omdb.OMDBResult `json:"details,omitempty"` // new field to hold OMDB data
}

type Torrent struct {
	URL          string `json:"url"`
	Hash         string `json:"hash"`
	Type         string `json:"type"`
	Seeds        int    `json:"seeds"`
	Peers        int    `json:"peers"`
	Size         string `json:"size"`
	SizeBytes    int64  `json:"size_bytes"`
	DateUploaded string `json:"date_uploaded"`
}

func SearchMovies(params Params) (*APIResponse, error) {
	searchPath := path + "?" + params.ToQueryParams()
	resp, err := http.Get(searchPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	return &apiResp, nil
}
