package internal

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/gin-gonic/gin"
)

var defaultTrackers = []string{
	"udp://open.stealth.si:80/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://tracker.coppersurfer.tk:6969/announce",
	"udp://exodus.desync.com:6969/announce",
	"udp://tracker.leechers-paradise.org:6969/announce",
}

var manager *TorrentManager

func init() {
	m, err := NewTorrentManager("./downloads") // put in a folder
	if err != nil {
		log.Fatal(err)
	}
	manager = m
}

func BuildMagnetLink(infoHash string, name string) string {
	var sb strings.Builder

	sb.WriteString("magnet:?xt=urn:btih:")
	sb.WriteString(strings.ToLower(infoHash))

	if name != "" {
		sb.WriteString("&dn=")
		sb.WriteString(url.QueryEscape(name))
	}

	for _, tr := range defaultTrackers {
		sb.WriteString("&tr=")
		sb.WriteString(url.QueryEscape(tr))
	}

	return sb.String()
}

type TorrentManager struct {
	client   *torrent.Client
	torrents map[string]*torrent.Torrent
	mu       sync.Mutex
}

func NewTorrentManager(downloadDir string) (*TorrentManager, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = downloadDir
	cfg.DisableIPv6 = true
	cfg.NoUpload = true
	cfg.Seed = false

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	log.Println("[INIT] Torrent client created")
	log.Println("[INIT] Download dir:", downloadDir)

	return &TorrentManager{
		client:   client,
		torrents: make(map[string]*torrent.Torrent),
	}, nil
}

func (m *TorrentManager) GetOrCreate(id string, magnet string) (*torrent.Torrent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t, ok := m.torrents[id]; ok {
		return t, nil
	}

	t, err := m.client.AddMagnet(magnet)
	if err != nil {
		return nil, err
	}

	<-t.GotInfo()

	// Pick largest file (movie)
	var file *torrent.File
	for _, f := range t.Files() {
		if file == nil || f.Length() > file.Length() {
			file = f
		}
	}

	file.SetPriority(torrent.PiecePriorityNow)
	file.Download()

	m.torrents[id] = t
	return t, nil
}

func Stream(c *gin.Context) {
	start := time.Now()

	log.Printf("hey")

	hash := c.Query("hash")
	log.Printf("hash : %s", hash)
	if hash == "" {
		c.String(http.StatusBadRequest, "missing id")
		return
	}

	log.Printf(
		"[HTTP] Client connected id=%s range=%s from=%s",
		hash,
		c.GetHeader("Range"),
		c.ClientIP(),
	)

	magnet := BuildMagnetLink(hash, "")
	t, err := manager.GetOrCreate(hash, magnet)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Pick the largest file
	var file *torrent.File
	for _, f := range t.Files() {
		if file == nil || f.Length() > file.Length() {
			file = f
		}
	}

	if file == nil {
		c.String(http.StatusNotFound, "no files found in torrent")
		return
	}

	log.Printf(
		"[STREAM] Start streaming hash=%s file=%s",
		hash,
		file.DisplayPath(),
	)

	reader := file.NewReader()
	defer func() {
		reader.Close()
		log.Printf(
			"[STREAM] Client disconnected hash=%s duration=%s",
			hash,
			time.Since(start),
		)
	}()

	// Streaming headers
	c.Header("Content-Type", "video/mp4")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Accept-Ranges", "bytes")

	// Make sure Gin doesn't buffer
	c.Writer.Flush()

	// FFmpeg command
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-movflags", "frag_keyframe+empty_moov+faststart",
		"-f", "mp4",
		"pipe:1",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := cmd.Start(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Torrent → FFmpeg stdin
	go func() {
		defer stdin.Close()
		if _, err := io.Copy(stdin, reader); err != nil {
			log.Printf("[ERROR] Copy to FFmpeg stdin: %v", err)
		}
	}()

	// FFmpeg stdout → client
	if _, err := io.Copy(c.Writer, stdout); err != nil {
		log.Printf("[ERROR] Copy FFmpeg stdout to client: %v", err)
	}

	cmd.Wait()
}
