package file

import (
	"encoding/json"
	"fmt"

	"115Quick_server/internal/service/auth"
)

type Service struct {
	client *auth.Client
}

func NewService(client *auth.Client) *Service {
	return &Service{client: client}
}

type FileInfo struct {
	FID       string `json:"fid"`
	PID       string `json:"pid"`
	FC        string `json:"fc"`
	FN        string `json:"fn"`
	FS        int64  `json:"fs"`
	SHA1      string `json:"sha1"`
	PC        string `json:"pc"`
	ICO       string `json:"ico"`
	IsV       int    `json:"isv"`
	Def       int    `json:"def"`
	Def2      int    `json:"def2"`
	PlayLong  int    `json:"play_long"`
	Thumb     string `json:"thumb"`
	UPT       int64  `json:"upt"`
	FTa       string `json:"fta"`
	IsM       string `json:"ism"`
}

type FileListResponse struct {
	Count   int        `json:"count"`
	Offset  int        `json:"offset"`
	Limit   int        `json:"limit"`
	CID     string     `json:"cid"`
	Data    []FileInfo `json:"data"`
	Path    []struct {
		Name string `json:"name"`
		CID  string `json:"cid"`
	} `json:"path"`
}

type FolderInfo struct {
	FileID      string `json:"file_id"`
	FileName    string `json:"file_name"`
	Size        string `json:"size"`
	SizeByte    int64  `json:"size_byte"`
	Count       string `json:"count"`
	FolderCount string `json:"folder_count"`
	PickCode    string `json:"pick_code"`
	SHA1        string `json:"sha1"`
	Ptime       string `json:"ptime"`
	Utime       string `json:"utime"`
	Paths       []struct {
		FileID   string `json:"file_id"`
		FileName string `json:"file_name"`
	} `json:"paths"`
}

func (s *Service) GetFileList(cid string, offset, limit int, fileType int) (*FileListResponse, error) {
	params := map[string]string{
		"offset": fmt.Sprintf("%d", offset),
		"limit":  fmt.Sprintf("%d", limit),
	}
	if cid != "" {
		params["cid"] = cid
	}
	if fileType > 0 {
		params["type"] = fmt.Sprintf("%d", fileType)
	}

	resp, err := s.client.Get("/open/ufile/files", params)
	if err != nil {
		return nil, fmt.Errorf("get file list: %w", err)
	}

	var result FileListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse file list: %w", err)
	}

	return &result, nil
}

func (s *Service) GetFileInfo(fileID string) (*FolderInfo, error) {
	params := map[string]string{
		"file_id": fileID,
	}

	resp, err := s.client.Get("/open/folder/get_info", params)
	if err != nil {
		return nil, fmt.Errorf("get file info: %w", err)
	}

	var result FolderInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse file info: %w", err)
	}

	return &result, nil
}

func (s *Service) SearchFiles(keyword string, offset, limit int) (*FileListResponse, error) {
	params := map[string]string{
		"search_value": keyword,
		"offset":       fmt.Sprintf("%d", offset),
		"limit":        fmt.Sprintf("%d", limit),
	}

	resp, err := s.client.Get("/open/ufile/search", params)
	if err != nil {
		return nil, fmt.Errorf("search files: %w", err)
	}

	var result FileListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse search results: %w", err)
	}

	return &result, nil
}

func (s *Service) CreateFolder(pid, name string) (string, error) {
	params := map[string]string{
		"pid":       pid,
		"file_name": name,
	}

	resp, err := s.client.Post("/open/folder/add", params, nil)
	if err != nil {
		return "", fmt.Errorf("create folder: %w", err)
	}

	var result struct {
		FileID string `json:"file_id"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", fmt.Errorf("parse create folder response: %w", err)
	}

	return result.FileID, nil
}

func (s *Service) DeleteFile(fileIDs string) error {
	params := map[string]string{
		"file_ids": fileIDs,
	}

	_, err := s.client.Post("/open/ufile/delete", params, nil)
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

func (s *Service) RenameFile(fileID, newName string) error {
	params := map[string]string{
		"file_id":   fileID,
		"file_name": newName,
	}

	_, err := s.client.Post("/open/ufile/update", params, nil)
	if err != nil {
		return fmt.Errorf("rename file: %w", err)
	}

	return nil
}

func (s *Service) MoveFile(fileIDs, toCID string) error {
	params := map[string]string{
		"file_ids": fileIDs,
		"to_cid":   toCID,
	}

	_, err := s.client.Post("/open/ufile/move", params, nil)
	if err != nil {
		return fmt.Errorf("move file: %w", err)
	}

	return nil
}

func (s *Service) GetDownloadURL(pickCode string) (string, error) {
	params := map[string]string{
		"pick_code": pickCode,
	}

	resp, err := s.client.Post("/open/ufile/downurl", params, nil)
	if err != nil {
		return "", fmt.Errorf("get download url: %w", err)
	}

	var result map[string]struct {
		URL struct {
			URL string `json:"url"`
		} `json:"url"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", fmt.Errorf("parse download url: %w", err)
	}

	for _, v := range result {
		if v.URL.URL != "" {
			return v.URL.URL, nil
		}
	}

	return "", fmt.Errorf("no download url found")
}
