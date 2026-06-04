package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"115Quick_server/internal/service/auth"
)

type Engine struct {
	tokenMgr     *auth.Manager
	downloadPath string
	progress     sync.Map
	active       sync.Map
	httpClient   *http.Client
}

type Progress struct {
	FileName    string  `json:"file_name"`
	FileSize    int64   `json:"file_size"`
	Downloaded  int64   `json:"downloaded"`
	Speed       float64 `json:"speed"`
	Percent     float64 `json:"percent"`
	Status      string  `json:"status"`
	Error       string  `json:"error,omitempty"`
	StartTime   int64   `json:"start_time"`
}

func NewEngine(tokenMgr *auth.Manager, downloadPath string) *Engine {
	return &Engine{
		tokenMgr:     tokenMgr,
		downloadPath: downloadPath,
		httpClient: &http.Client{
			Timeout: 30 * time.Minute,
		},
	}
}

func (e *Engine) Download(ctx context.Context, url, fileName string, fileSize int64) error {
	key := fileName

	if _, exists := e.active.LoadOrStore(key, true); exists {
		return fmt.Errorf("download already in progress: %s", fileName)
	}
	defer e.active.Delete(key)

	if err := os.MkdirAll(e.downloadPath, 0755); err != nil {
		return fmt.Errorf("create download directory: %w", err)
	}

	filePath := filepath.Join(e.downloadPath, fileName)
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	token, err := e.tokenMgr.GetAccessToken()
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	tmpPath := filePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	prog := &Progress{
		FileName:  fileName,
		FileSize:  fileSize,
		Status:    "downloading",
		StartTime: time.Now().Unix(),
	}
	e.progress.Store(key, prog)

	var downloaded int64
	var lastDownloaded int64
	var speed float64

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				current := atomic.LoadInt64(&downloaded)
				speed = float64(current-lastDownloaded)
				lastDownloaded = current
				prog.Downloaded = current
				prog.Speed = speed
				if fileSize > 0 {
					prog.Percent = float64(current) / float64(fileSize) * 100
				}
			case <-done:
				return
			}
		}
	}()

	writer := &progressWriter{writer: f, downloaded: &downloaded}
	_, err = io.Copy(writer, resp.Body)
	close(done)

	if err != nil {
		os.Remove(tmpPath)
		prog.Status = "failed"
		prog.Error = err.Error()
		return fmt.Errorf("download file: %w", err)
	}

	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("rename file: %w", err)
	}

	prog.Status = "completed"
	prog.Percent = 100
	prog.Downloaded = fileSize

	return nil
}

func (e *Engine) GetProgress(fileName string) *Progress {
	if v, ok := e.progress.Load(fileName); ok {
		return v.(*Progress)
	}
	return nil
}

func (e *Engine) GetAllProgress() map[string]*Progress {
	result := make(map[string]*Progress)
	e.progress.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(*Progress)
		return true
	})
	return result
}

func (e *Engine) RemoveProgress(fileName string) {
	e.progress.Delete(fileName)
}

type progressWriter struct {
	writer     io.Writer
	downloaded *int64
}

func (w *progressWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	atomic.AddInt64(w.downloaded, int64(n))
	return
}
