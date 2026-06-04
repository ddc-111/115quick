package video

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"115Quick_server/internal/service/auth"
)

type Service struct {
	client   *auth.Client
	tokenMgr *auth.Manager
}

func NewService(client *auth.Client, tokenMgr *auth.Manager) *Service {
	return &Service{
		client:   client,
		tokenMgr: tokenMgr,
	}
}

type PlayInfo struct {
	FileID      string `json:"file_id"`
	ParentID    string `json:"parent_id"`
	FileName    string `json:"file_name"`
	FileSize    string `json:"file_size"`
	FileSHA1    string `json:"file_sha1"`
	FileType    string `json:"file_type"`
	PlayLong    string `json:"play_long"`
	IsPrivate   string `json:"is_private"`
	UserDef     int    `json:"user_def"`
	VideoURL    []struct {
		URL        string `json:"url"`
		Height     int    `json:"height"`
		Width      int    `json:"width"`
		Definition int    `json:"definition"`
		Title      string `json:"title"`
	} `json:"video_url"`
	DefinitionListNew []int `json:"definition_list_new"`
	MultitrackList    []struct {
		Title      string `json:"title"`
		IsSelected string `json:"is_selected"`
	} `json:"multitrack_list"`
}

type SubtitleInfo struct {
	Autoload []struct {
		SID      string `json:"sid"`
		Language string `json:"language"`
		Title    string `json:"title"`
		URL      string `json:"url"`
		Type     string `json:"type"`
	} `json:"autoload"`
	List []struct {
		SID      string `json:"sid"`
		Language string `json:"language"`
		Title    string `json:"title"`
		URL      string `json:"url"`
		Type     string `json:"type"`
		SHA1     string `json:"sha1"`
		FileID   string `json:"file_id"`
		FileName string `json:"file_name"`
		PickCode string `json:"pick_code"`
	} `json:"list"`
}

type PlayHistory struct {
	AddTime  int64  `json:"add_time"`
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	PickCode string `json:"pick_code"`
	Time     string `json:"time"`
}

func (s *Service) GetPlayInfo(pickCode string) (*PlayInfo, error) {
	params := map[string]string{
		"pick_code": pickCode,
	}

	resp, err := s.client.Get("/open/video/play", params)
	if err != nil {
		return nil, fmt.Errorf("get play info: %w", err)
	}

	var result PlayInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse play info: %w", err)
	}

	return &result, nil
}

func (s *Service) GetSubtitles(pickCode string) (*SubtitleInfo, error) {
	params := map[string]string{
		"pick_code": pickCode,
	}

	resp, err := s.client.Get("/open/video/subtitle", params)
	if err != nil {
		return nil, fmt.Errorf("get subtitles: %w", err)
	}

	var result SubtitleInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse subtitles: %w", err)
	}

	return &result, nil
}

func (s *Service) GetPlayHistory(pickCode string) (*PlayHistory, error) {
	params := map[string]string{
		"pick_code": pickCode,
	}

	resp, err := s.client.Get("/open/video/history", params)
	if err != nil {
		return nil, fmt.Errorf("get play history: %w", err)
	}

	var result PlayHistory
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse play history: %w", err)
	}

	return &result, nil
}

func (s *Service) SavePlayHistory(pickCode string, timeSec int, watchEnd bool) error {
	params := map[string]string{
		"pick_code": pickCode,
		"time":      fmt.Sprintf("%d", timeSec),
	}
	if watchEnd {
		params["watch_end"] = "1"
	}

	_, err := s.client.Post("/open/video/history", params, nil)
	if err != nil {
		return fmt.Errorf("save play history: %w", err)
	}

	return nil
}

func (s *Service) SubmitTranscode(pickCode, op string) error {
	params := map[string]string{
		"pick_code": pickCode,
		"op":        op,
	}

	_, err := s.client.Post("/open/video/video_push", params, nil)
	if err != nil {
		return fmt.Errorf("submit transcode: %w", err)
	}

	return nil
}

func (s *Service) ProxyVideoStream(pickCode, definition string, w http.ResponseWriter, r *http.Request) error {
	playInfo, err := s.GetPlayInfo(pickCode)
	if err != nil {
		return fmt.Errorf("get play info: %w", err)
	}

	var targetURL string
	if definition != "" {
		for _, v := range playInfo.VideoURL {
			if fmt.Sprintf("%d", v.Definition) == definition {
				targetURL = v.URL
				break
			}
		}
	}

	if targetURL == "" && len(playInfo.VideoURL) > 0 {
		targetURL = playInfo.VideoURL[0].URL
	}

	if targetURL == "" {
		return fmt.Errorf("no video url available")
	}

	token, err := s.tokenMgr.GetAccessToken()
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Minute}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("create proxy request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	for key, values := range r.Header {
		for _, value := range values {
			if key != "Authorization" && key != "Host" {
				req.Header.Add(key, value)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("proxy request: %w", err)
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	return err
}
