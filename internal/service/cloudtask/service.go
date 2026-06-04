package cloudtask

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

type Task struct {
	InfoHash     string `json:"info_hash"`
	AddTime      int64  `json:"add_time"`
	PercentDone  int    `json:"percentDone"`
	Size         int64  `json:"size"`
	Name         string `json:"name"`
	LastUpdate   int64  `json:"last_update"`
	FileID       string `json:"file_id"`
	DeleteFileID string `json:"delete_file_id"`
	Status       int    `json:"status"`
	URL          string `json:"url"`
	WpPathID     string `json:"wp_path_id"`
	Def2         int    `json:"def2"`
	PlayLong     int    `json:"play_long"`
}

type TaskListResponse struct {
	Page      int    `json:"page"`
	PageCount int    `json:"page_count"`
	Count     int    `json:"count"`
	Tasks     []Task `json:"tasks"`
}

type QuotaInfo struct {
	Count   int `json:"count"`
	Surplus int `json:"surplus"`
	Used    int `json:"used"`
	Package []struct {
		Surplus int    `json:"surplus"`
		Used    int    `json:"used"`
		Count   int    `json:"count"`
		Name    string `json:"name"`
	} `json:"package"`
}

type AddTaskResult struct {
	State   bool   `json:"state"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	InfoHash string `json:"info_hash"`
	URL      string `json:"url"`
}

func (s *Service) GetTaskList(page int) (*TaskListResponse, error) {
	params := map[string]string{}
	if page > 1 {
		params["page"] = fmt.Sprintf("%d", page)
	}

	resp, err := s.client.Get("/open/offline/get_task_list", params)
	if err != nil {
		return nil, fmt.Errorf("get task list: %w", err)
	}

	var result TaskListResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse task list: %w", err)
	}

	return &result, nil
}

func (s *Service) AddTaskURLs(urls string, wpPathID string) ([]AddTaskResult, error) {
	params := map[string]string{
		"urls": urls,
	}
	if wpPathID != "" {
		params["wp_path_id"] = wpPathID
	}

	resp, err := s.client.Post("/open/offline/add_task_urls", params, nil)
	if err != nil {
		return nil, fmt.Errorf("add task urls: %w", err)
	}

	var result []AddTaskResult
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse add task result: %w", err)
	}

	return result, nil
}

func (s *Service) AddTaskBT(infoHash, wanted, savePath, torrentSHA1, pickCode, wpPathID string) error {
	params := map[string]string{
		"info_hash":    infoHash,
		"wanted":       wanted,
		"save_path":    savePath,
		"torrent_sha1": torrentSHA1,
		"pick_code":    pickCode,
	}
	if wpPathID != "" {
		params["wp_path_id"] = wpPathID
	}

	_, err := s.client.Post("/open/offline/add_task_bt", params, nil)
	if err != nil {
		return fmt.Errorf("add bt task: %w", err)
	}

	return nil
}

func (s *Service) DeleteTask(infoHash string, deleteSourceFile bool) error {
	params := map[string]string{
		"info_hash": infoHash,
	}
	if deleteSourceFile {
		params["del_source_file"] = "1"
	}

	_, err := s.client.Post("/open/offline/del_task", params, nil)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}

func (s *Service) ClearTasks(flag int) error {
	params := map[string]string{
		"flag": fmt.Sprintf("%d", flag),
	}

	_, err := s.client.Post("/open/offline/clear_task", params, nil)
	if err != nil {
		return fmt.Errorf("clear tasks: %w", err)
	}

	return nil
}

func (s *Service) GetQuotaInfo() (*QuotaInfo, error) {
	resp, err := s.client.Get("/open/offline/get_quota_info", nil)
	if err != nil {
		return nil, fmt.Errorf("get quota info: %w", err)
	}

	var result QuotaInfo
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("parse quota info: %w", err)
	}

	return &result, nil
}

func (s *Service) ParseTorrent(torrentSHA1, pickCode string) (interface{}, error) {
	params := map[string]string{
		"torrent_sha1": torrentSHA1,
		"pick_code":    pickCode,
	}

	resp, err := s.client.Post("/open/offline/torrent", params, nil)
	if err != nil {
		return nil, fmt.Errorf("parse torrent: %w", err)
	}

	return resp.Data, nil
}
