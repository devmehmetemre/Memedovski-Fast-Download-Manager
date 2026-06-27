package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type chunkInfo struct {
	start int64
	end   int64
}

type writeReq struct {
	offset int64
	data   []byte
}

var (
	sharedTransport = &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     200,
		IdleConnTimeout:     90 * time.Second,
		ForceAttemptHTTP2:   true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		WriteBufferSize:       512 * 1024,
		ReadBufferSize:        512 * 1024,
	}
	sharedClient = &http.Client{
		Transport: sharedTransport,
		Timeout:   0,
	}

	mimeExtensions = map[string]string{
		"application/zip":                 ".zip",
		"application/x-zip-compressed":    ".zip",
		"application/x-rar-compressed":    ".rar",
		"application/x-7z-compressed":     ".7z",
		"application/gzip":                ".gz",
		"application/x-gzip":              ".gz",
		"application/x-tar":              ".tar",
		"application/x-bzip2":            ".bz2",
		"application/x-xz":               ".xz",
		"application/pdf":                ".pdf",
		"application/msword":             ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"application/vnd.ms-powerpoint": ".ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
		"application/json":               ".json",
		"application/xml":                ".xml",
		"application/octet-stream":       ".bin",
		"application/x-msdownload":       ".exe",
		"application/x-msi":              ".msi",
		"application/x-iso9660-image":    ".iso",
		"application/x-shockwave-flash": ".swf",
		"application/java-archive":       ".jar",
		"application/x-www-form-urlencoded": ".dat",
		"text/plain":                     ".txt",
		"text/html":                      ".html",
		"text/css":                       ".css",
		"text/javascript":                ".js",
		"text/csv":                       ".csv",
		"text/xml":                       ".xml",
		"text/markdown":                  ".md",
		"image/jpeg":                     ".jpg",
		"image/png":                      ".png",
		"image/gif":                      ".gif",
		"image/webp":                     ".webp",
		"image/bmp":                      ".bmp",
		"image/svg+xml":                  ".svg",
		"image/x-icon":                   ".ico",
		"audio/mpeg":                     ".mp3",
		"audio/wav":                      ".wav",
		"audio/ogg":                      ".ogg",
		"audio/flac":                     ".flac",
		"audio/aac":                      ".aac",
		"video/mp4":                      ".mp4",
		"video/x-msvideo":               ".avi",
		"video/x-matroska":              ".mkv",
		"video/quicktime":               ".mov",
		"video/webm":                     ".webm",
		"video/mpeg":                     ".mpeg",
		"video/x-ms-wmv":                ".wmv",
		"application/x-bittorrent":       ".torrent",
	}
)

func writeWorker(f *os.File, writes <-chan writeReq, dld *atomic.Int64, cancel <-chan struct{}) {
	for wr := range writes {
		select {
		case <-cancel:
			return
		default:
		}
		if _, err := f.WriteAt(wr.data, wr.offset); err != nil {
			return
		}
		if dld != nil {
			dld.Add(int64(len(wr.data)))
		}
	}
}

func downloadChunk(urlStr string, ci chunkInfo, writes chan<- writeReq, cancel <-chan struct{}) error {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", ci.start, ci.end))

	var resp *http.Response
	for attempt := 0; attempt < 5; attempt++ {
		if cancel != nil {
			select {
			case <-cancel:
				return nil
			default:
			}
		}
		resp, err = sharedClient.Do(req)
		if err == nil {
			if resp.StatusCode == http.StatusPartialContent || resp.StatusCode == http.StatusOK {
				break
			}
			resp.Body.Close()
		}
		if attempt < 4 {
			time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
		}
	}
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	offset := ci.start
	buf := make([]byte, 1024*1024)
	for {
		if cancel != nil {
			select {
			case <-cancel:
				return nil
			default:
			}
		}
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			writes <- writeReq{offset: offset, data: data}
			offset += int64(n)
		}
		if rerr == io.EOF {
			return nil
		}
		if rerr != nil {
			return fmt.Errorf("read: %w", rerr)
		}
	}
}

func fetchFileInfo(urlStr string) (fileSize int64, rangesSupported bool, filename string, err error) {
	resp, err := sharedClient.Head(urlStr)
	if err != nil {
		return 0, false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return 0, false, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	acceptRanges := resp.Header.Get("Accept-Ranges")
	rangesSupported = strings.EqualFold(acceptRanges, "bytes")

	sizeStr := resp.Header.Get("Content-Length")
	if sizeStr != "" {
		fileSize, _ = strconv.ParseInt(sizeStr, 10, 64)
	}

	filename = resolveFilename(urlStr, resp.Header)
	return
}

func resolveFilename(rawURL string, headers http.Header) string {
	cd := headers.Get("Content-Disposition")
	if cd != "" {
		for _, part := range strings.Split(cd, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "filename=") || strings.HasPrefix(part, "filename*=UTF-8''") {
				if strings.HasPrefix(part, "filename*=UTF-8''") {
					part = strings.TrimPrefix(part, "filename*=UTF-8''")
				} else {
					part = strings.TrimPrefix(part, "filename=")
				}
				part = strings.Trim(part, "\"")
				if part != "" {
					return part
				}
			}
		}
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Sprintf("download_%d", time.Now().Unix())
	}

	base := path.Base(u.Path)
	if base == "" || base == "." || base == "/" || base == "download" {
		base = fmt.Sprintf("download_%d", time.Now().Unix())
	}

	if strings.Contains(base, ".") {
		return base
	}

	ct := headers.Get("Content-Type")
	ct = strings.Split(ct, ";")[0]
	ct = strings.TrimSpace(ct)

	if ext, ok := mimeExtensions[ct]; ok {
		return base + ext
	}

	if strings.HasPrefix(ct, "text/") {
		return base + ".txt"
	}
	if strings.HasPrefix(ct, "image/") {
		return base + ".img"
	}
	if strings.HasPrefix(ct, "video/") {
		return base + ".vid"
	}
	if strings.HasPrefix(ct, "audio/") {
		return base + ".aud"
	}
	if strings.HasPrefix(ct, "application/") {
		return base + ".bin"
	}

	return base
}

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatSpeed(b int64, d time.Duration) string {
	secs := d.Seconds()
	if secs <= 0 {
		return ""
	}
	return formatSize(int64(float64(b)/secs)) + "/s"
}

func formatDuration(d time.Duration) string {
	if d.Hours() >= 1 {
		return fmt.Sprintf("%dh%dm%ds", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
	}
	if d.Minutes() >= 1 {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func formatETA(downloaded, total int64, elapsed time.Duration) string {
	if downloaded <= 0 || total <= 0 {
		return ""
	}
	remaining := float64(total-downloaded) / float64(downloaded) * elapsed.Seconds()
	if remaining <= 0 {
		return ""
	}
	return "ETA: " + formatDuration(time.Duration(remaining)*time.Second/time.Nanosecond)
}

func progressBar(pct float64, width int) string {
	filled := int(pct * float64(width) / 100.0)
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("=", filled) + strings.Repeat(" ", width-filled) + "]"
}
