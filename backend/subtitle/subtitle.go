package subtitle

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strconv"
	// "os"
    "github.com/gin-gonic/gin"
)

type Subtitle struct {
    ID           int    `json:"id"`
    FileID       int    `json:"file_id"`
    Filename     string `json:"filename"`
    Language     string `json:"language"`
    DownloadLink string `json:"download_link,omitempty"`
}

type searchResponse struct {
    Data []Subtitle `json:"data"`
}

type downloadRequest struct {
    FileID        int    `json:"file_id"`
    SubFormat     string `json:"sub_format,omitempty"`
    FileName      string `json:"file_name,omitempty"`
    ForceDownload bool   `json:"force_download,omitempty"`
}

type downloadResponse struct {
    Link         string `json:"link"`
    FileName     string `json:"file_name"`
    Requests     int    `json:"requests"`
    Allowed      int    `json:"allowed"`
    ResetTimeUTC string `json:"reset_time_utc"`
    Message      string `json:"message"`
}


var apikey string = "7bAYMknSYdSZ9TgdYRwI0SLQ4jl2io8i"

func openSubtitlesHeaders(req *http.Request) {
    req.Header.Set("Api-Key", apikey)
    req.Header.Set("User-Agent", "MyGoApp v1.0")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
}

func isValidInt(s string) bool {
    _, err := strconv.Atoi(s)
    return err == nil
}

func SearchSubtitlesHandler(c *gin.Context) {
    query := c.Query("query")
    if query == "" {
        c.JSON(400, gin.H{"error": "query parameter is required"})
        return
    }

    params := url.Values{}
    params.Set("query", query)
    lang := c.DefaultQuery("lang", "en")
    params.Set("languages", lang)

    season := c.Query("season")
    episode := c.Query("episode")
    subType := c.Query("type")
    page := c.DefaultQuery("page", "1")

    if season != "" && isValidInt(season) {
        params.Set("season_number", season)
    }
    if episode != "" && isValidInt(episode) {
        params.Set("episode_number", episode)
    }
    if subType == "movie" || subType == "episode" {
        params.Set("type", subType)
    }
    if isValidInt(page) {
        params.Set("page", page)
    }

    apiURL := "https://api.opensubtitles.com/api/v1/subtitles?" + params.Encode()

    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    openSubtitlesHeaders(req)

    resp, err := (&http.Client{}).Do(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result interface{}
    if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(500, gin.H{"error": "failed to parse response: " + err.Error()})
        return 
    }
    c.JSON(200, result)
}

func GetDownloadLinkHandler(c *gin.Context) {
    var body downloadRequest
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(400, gin.H{"error": "invalid body: " + err.Error()})
        return
    }
    if body.FileID == 0 {
        c.JSON(400, gin.H{"error": "file_id is required"})
        return
    }
    if body.SubFormat == "" {
        body.SubFormat = "srt"
    }

    payload, _ := json.Marshal(body)
    req, err := http.NewRequest("POST", "https://api.opensubtitles.com/api/v1/download", bytes.NewBuffer(payload))
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    openSubtitlesHeaders(req)

    resp, err := (&http.Client{}).Do(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        var errResp interface{}
        json.Unmarshal(respBody, &errResp)
        c.JSON(resp.StatusCode, errResp)
        return
    }

    var dlResp downloadResponse
    if err := json.Unmarshal(respBody, &dlResp); err != nil {
        c.JSON(500, gin.H{"error": "failed to parse download response"})
        return
    }

    c.JSON(200, gin.H{
        "link":            dlResp.Link,
        "file_name":       dlResp.FileName,
        "quota_remaining": dlResp.Requests,
        "quota_allowed":   dlResp.Allowed,
        "reset_time_utc":  dlResp.ResetTimeUTC,
    })
}

func DownloadSubtitleFileHandler(c *gin.Context) {
    link := c.Query("link")
    if link == "" {
        c.JSON(400, gin.H{"error": "link parameter is required"})
        return
    }

    parsed, err := url.ParseRequestURI(link)
    if err != nil || (parsed.Host != "dl.opensubtitles.com" &&
        parsed.Host != "www.opensubtitles.com") {
        c.JSON(400, gin.H{"error": "invalid or untrusted download link"})
        return
    }

    req, err := http.NewRequest("GET", link, nil)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    req.Header.Set("User-Agent", "MyGoApp v1.0")

    resp, err := (&http.Client{}).Do(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("upstream returned %d", resp.StatusCode)})
        return
    }

    fileName := c.DefaultQuery("file_name", "subtitle.srt")
    c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
    c.Header("Content-Type", "text/plain; charset=utf-8")

    io.Copy(c.Writer, resp.Body)
}