package subtitle

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "regexp"
    "sort"
    "strconv"
    "github.com/gin-gonic/gin"
)

type Subtitle struct {
    ID           int    `json:id`
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

var sePattern = regexp.MustCompile(`(?i)S(\d+)E(\d+)`)

func parseSeasonEpisode(filename string) (int, int) {
    matches := sePattern.FindStringSubmatch(filename)
    if len(matches) == 3 {
        s, _ := strconv.Atoi(matches[1])
        e, _ := strconv.Atoi(matches[2])
        return s, e
    }
    return 0, 0
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

    var raw struct {
        Data []struct {
            Attributes struct {
                FeatureDetails struct {
                    MovieName     string `json:"movie_name"`
                    FeatureType   string `json:"feature_type"`
                    SeasonNumber  int    `json:"season_number"`
                    EpisodeNumber int    `json:"episode_number"`
                } `json:"feature_details"`
                Files []struct {
                    FileID   int    `json:"file_id"`
                    FileName string `json:"file_name"`
                } `json:"files"`
                Language string `json:"language"`
            } `json:"attributes"`
        } `json:"data"`
    }

    if err := json.Unmarshal(body, &raw); err != nil {
        c.JSON(500, gin.H{"error": "failed to parse response: " + err.Error()})
        return
    }

    if len(raw.Data) == 0 {
        c.JSON(404, gin.H{"error": "no subtitles found"})
        return
    }

    movieName := query
    if raw.Data[0].Attributes.FeatureDetails.MovieName != "" {
        movieName = raw.Data[0].Attributes.FeatureDetails.MovieName
    }

    type subtitleResult struct {
        Season   int    `json:"season"`
        Episode  int    `json:"episode"`
        FileID   int    `json:"file_id"`
        Filename string `json:"filename"`
        Language string `json:"language"`
    }

    results := make([]subtitleResult, 0, len(raw.Data))
    seen := make(map[string]bool)
    for _, item := range raw.Data {
        for _, f := range item.Attributes.Files {
            s := item.Attributes.FeatureDetails.SeasonNumber
            e := item.Attributes.FeatureDetails.EpisodeNumber
            if s == 0 && e == 0 {
                s, e = parseSeasonEpisode(f.FileName)
            }
            key := strconv.Itoa(s) + "-" + strconv.Itoa(e)
            if seen[key] {
                continue
            }
            seen[key] = true
            results = append(results, subtitleResult{
                Season:   s,
                Episode:  e,
                FileID:   f.FileID,
                Filename: f.FileName,
                Language: item.Attributes.Language,
            })
        }
    }

    sort.Slice(results, func(i, j int) bool {
        if results[i].Season != results[j].Season {
            return results[i].Season < results[j].Season
        }
        return results[i].Episode < results[j].Episode
    })

    c.JSON(200, gin.H{
        "movie_name": movieName,
        "results":    results,
    })
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