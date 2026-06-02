package svc

import (
	"115Quick_server/internal/types"
	"115Quick_server/internal/utils/request115"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileInfo struct {
	InfoHash string    `json:"infoHash"`
	URL      string    `json:"url"`
	AddTime  time.Time `json:"addTime"`
}

type Auth115Manager struct {
	mu            sync.RWMutex
	authInfo      *types.Auth115Info
	lastRefreshAt time.Time
	svcCtx        *ServiceContext
	downloadPath  string

	fileInfos        []FileInfo
	fileMutex        sync.Mutex
	lastPage         int
	DownloadInfoChan chan *DownloadInfo
	CloudTask        AllTasks
	running          int32
}

type AllTasks []struct {
	InfoHash     string  `json:"info_hash"`
	AddTime      int64   `json:"add_time"`
	PercentDone  float64 `json:"percentDone"`
	Size         int64   `json:"size"`
	Peers        int     `json:"peers"`
	RateDownload int     `json:"rateDownload"`
	Name         string  `json:"name"`
	LastUpdate   int64   `json:"last_update"`
	LeftTime     int     `json:"left_time"`
	FileID       string  `json:"file_id"`
	DeleteFileID string  `json:"delete_file_id"`
	Move         int     `json:"move"`
	Status       int     `json:"status"`
	Err          int     `json:"err"`
	URL          string  `json:"url"`
	DelPath      string  `json:"del_path"`
	WpPathID     string  `json:"wp_path_id"`
	Def2         int     `json:"def2"`
	PlayLong     *int    `json:"play_long"`
	CanAppeal    int     `json:"can_appeal"`
}

type DownloadInfo struct {
	FolderName  string
	fileID      string
	fileName    string
	fileSize    int64
	downloadURL string
}

func NewAuth115Manager(svcCtx *ServiceContext) *Auth115Manager {
	m := &Auth115Manager{
		svcCtx:           svcCtx,
		fileInfos:        make([]FileInfo, 0),
		lastPage:         1,
		DownloadInfoChan: make(chan *DownloadInfo, 1000000),
		running:          0,
	}

	accessToken := svcCtx.Config.Auth115.AccessToken
	refreshToken := svcCtx.Config.Auth115.RefreshToken
	if accessToken != "" && refreshToken != "" {
		m.authInfo = &types.Auth115Info{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    7200,
		}
		m.lastRefreshAt = time.Now()
		request115.SetAccessToken(accessToken)
	}

	m.downloadPath = svcCtx.Config.Auth115.DownloadPath

	m.loadFileInfos()

	m.startFilePolling()
	m.startFileInfoCheck()
	m.startDownload()
	m.PollFiles()
	return m
}

func (m *Auth115Manager) startDownload() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for info := range m.DownloadInfoChan {
			<-ticker.C
			err := m.downloadFile(info.FolderName, info.fileID, info.fileName, info.fileSize, info.downloadURL)
			if err != nil {
				logx.Error(err)
				go func() {
					time.Sleep(time.Hour * 1)
					m.DownloadInfoChan <- &DownloadInfo{
						FolderName:  info.FolderName,
						fileID:      info.fileID,
						fileName:    info.fileName,
						fileSize:    info.fileSize,
						downloadURL: info.downloadURL,
					}
				}()
			}
		}
	}()
}

func (m *Auth115Manager) GetAccessToken() (string, error) {
	m.mu.RLock()
	if m.authInfo == nil {
		m.mu.RUnlock()
		return "", fmt.Errorf("未设置115认证信息，请在配置文件中设置AccessToken和RefreshToken")
	}

	expiresAt := m.lastRefreshAt.Add(time.Duration(m.authInfo.ExpiresIn) * time.Second)
	if time.Now().Add(5 * time.Minute).Before(expiresAt) {
		token := m.authInfo.AccessToken
		m.mu.RUnlock()
		return token, nil
	}
	m.mu.RUnlock()

	return m.RefreshToken()
}

func (m *Auth115Manager) RefreshToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.authInfo != nil {
		expiresAt := m.lastRefreshAt.Add(time.Duration(m.authInfo.ExpiresIn) * time.Second)
		if time.Now().Add(5 * time.Minute).Before(expiresAt) {
			return m.authInfo.AccessToken, nil
		}
	}

	data := url.Values{}
	data.Set("refresh_token", m.authInfo.RefreshToken)

	apiURL := "https://passportapi.115.com/open/refreshToken"
	request, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建刷新请求失败: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("发送刷新请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("读取刷新响应失败: %v", err)
	}

	var result struct {
		State int `json:"state"`
		Code  int `json:"code"`
		Data  struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
		} `json:"data"`
		Error string `json:"error"`
		Errno int    `json:"errno"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析刷新响应失败: %v", err)
	}

	if result.State != 1 || result.Code != 0 || result.Errno != 0 {
		return "", fmt.Errorf("刷新令牌失败: %s", result.Error)
	}

	m.authInfo = &types.Auth115Info{
		AccessToken:  result.Data.AccessToken,
		RefreshToken: result.Data.RefreshToken,
		ExpiresIn:    result.Data.ExpiresIn,
	}
	m.lastRefreshAt = time.Now()

	request115.SetAccessToken(m.authInfo.AccessToken)

	logx.Infof("令牌刷新成功，新令牌已同步到request115")

	return m.authInfo.AccessToken, nil
}

func (m *Auth115Manager) GetAuthInfo() *types.Auth115Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.authInfo == nil {
		return nil
	}

	info := *m.authInfo
	return &info
}

func (m *Auth115Manager) AddFileInfo(infoHash string, url string) {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()

	fileInfo := FileInfo{
		InfoHash: infoHash,
		URL:      url,
		AddTime:  time.Now(),
	}
	m.fileInfos = append(m.fileInfos, fileInfo)

	m.saveFileInfos()
}

func (m *Auth115Manager) saveFileInfos() {
	if err := m.svcCtx.SaveFileInfos(m.fileInfos); err != nil {
		logx.Errorf("保存文件信息失败: %v", err)
	}
}

func (m *Auth115Manager) fileInfoCheck() bool {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()
	if len(m.fileInfos) == 0 {
		return false
	}
	for i := len(m.fileInfos) - 1; i >= 0; i-- {
		fileInfo := m.fileInfos[i]
		for _, task := range m.CloudTask {
			if task.URL == fileInfo.URL {
				if task.Status == 2 {
					logx.Infof("文件下载完成，文件ID: %s", task.FileID)
					m.fileInfos = append(m.fileInfos[:i], m.fileInfos[i+1:]...)
					m.saveFileInfos()

					fileDetail, err := m.GetFileInfo(task.FileID)

					if err != nil {
						logx.Errorf("获取文件详情失败: %v", err)
					} else {
						logx.Infof("文件提取码: %s", fileDetail.Data.PickCode)

						if fileDetail.Data.FileCategory == "0" {
							logx.Infof("检测到文件夹: %s，获取文件夹内视频文件提取码", fileDetail.Data.FileName)

							var FileTypeId int
							if m.svcCtx.Mode == 0 {
								FileTypeId = 4
							} else {
								FileTypeId = 0
							}
							folderFiles, err := m.GetFolderFiles(task.FileID, FileTypeId)

							if err != nil {
								logx.Errorf("获取文件夹内文件列表失败: %v", err)
							} else {
								if len(folderFiles.Data) > 0 {
									logx.Infof("文件夹 %s 内共有 %d 个视频文件", fileDetail.Data.FileName, len(folderFiles.Data))
									for i, file := range folderFiles.Data {
										logx.Infof("视频文件 %d: %s, 提取码: %s", i+1, file.FN, file.PC)
										downloadUrlResp, err := m.GetFileDownloadUrl(file.PC)
										if err != nil {
											logx.Errorf("获取文件下载地址失败: %v", err)
										} else {
											for fileID, fileInfo := range downloadUrlResp.Data {
												logx.Infof("文件ID: %s, 文件名: %s, 文件大小: %d, 下载地址: %s",
													fileID, fileInfo.FileName, fileInfo.FileSize, fileInfo.URL.URL)

												m.DownloadInfoChan <- &DownloadInfo{
													FolderName:  fileDetail.Data.FileName,
													fileID:      fileID,
													fileName:    fileInfo.FileName,
													fileSize:    fileInfo.FileSize,
													downloadURL: fileInfo.URL.URL,
												}
											}
										}
									}
								} else {
									logx.Infof("文件夹 %s 内没有视频文件", fileDetail.Data.FileName)
								}
							}
							return true
						}
					}
				}
				break
			}
		}
	}
	return false
}

func (m *Auth115Manager) loadFileInfos() {
	fileInfos := m.svcCtx.LoadFileInfos()
	if fileInfos != nil {
		m.fileInfos = fileInfos
	}
}

func (m *Auth115Manager) startFilePolling() {
	ticker := time.NewTicker(60 * time.Minute)
	go func() {
		for range ticker.C {
			m.PollFiles()
		}
	}()
}

func (m *Auth115Manager) DoOnceWithCooldownPoll() {
	if !atomic.CompareAndSwapInt32(&m.running, 0, 1) {
		fmt.Println("请求被忽略（5分钟冷却中）")
		return
	}

	go func() {
		m.PollFiles()
		time.Sleep(4 * time.Minute)
		atomic.StoreInt32(&m.running, 0)
	}()
}

func (m *Auth115Manager) startFileInfoCheck() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			if m.fileInfoCheck() {
				time.Sleep(60 * time.Second)
			}
		}
	}()
}

func (m *Auth115Manager) GetOfflineTaskList(page int) (*types.GetOfflineTaskListResp, error) {
	accessToken, err := m.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}

	params := url.Values{}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}

	apiURL := "https://proapi.115.com/open/offline/get_task_list"
	if len(params) > 0 {
		apiURL += "?" + params.Encode()
	}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result types.GetOfflineTaskListResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", body)
	}

	if !result.State {
		return nil, fmt.Errorf("获取云下载任务列表失败: %s", result.Message)
	}

	return &result, nil
}

func (m *Auth115Manager) PollFiles() {
	logx.Info("开始获取云下载任务")
	var page = 1
	m.CloudTask = nil
	for {
		resp, err := m.GetOfflineTaskList(page)
		if err != nil {
			logx.Errorf("获取云下载任务列表失败: %v", err)
			return
		}

		if resp.Data.Count == 0 {
			logx.Info("没有云下载任务")
			return
		}
		m.CloudTask = append(m.CloudTask, resp.Data.Tasks...)

		if page >= resp.Data.PageCount {
			break
		}

		page++
	}
}

func (m *Auth115Manager) GetFileInfo(fileID string) (*types.FileInfoResp, error) {
	accessToken, err := m.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}

	params := url.Values{}
	params.Set("file_id", fileID)

	apiURL := "https://proapi.115.com/open/folder/get_info?" + params.Encode()
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result types.FileInfoResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.State {
		return nil, fmt.Errorf("获取文件详情失败: %s", result.Message)
	}

	return &result, nil
}

func (m *Auth115Manager) GetFileDownloadUrl(pickCode string) (*types.GetFileDownloadUrlResp, error) {
	accessToken, err := m.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}

	data := url.Values{}
	data.Set("pick_code", pickCode)

	apiURL := "https://proapi.115.com/open/ufile/downurl"
	request, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result types.GetFileDownloadUrlResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.State {
		return nil, fmt.Errorf("获取文件下载地址失败: %s", result.Message)
	}

	return &result, nil
}

func (m *Auth115Manager) downloadFile(FolderName string, fileID string, fileName string, fileSize int64, downloadURL string) error {
	dataDir := m.downloadPath + "/" + FolderName
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("创建data目录失败: %v", err)
		}
	}

	filePath := filepath.Join(dataDir, fileName)

	if _, err := os.Stat(filePath); err == nil {
		logx.Infof("文件已存在，跳过下载: %s", filePath)
		return nil
	}

	request, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %v", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("发送下载请求失败: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("下载请求失败，状态码: %d", response.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	progress := &Progress{
		Total:     fileSize,
		Current:   0,
		StartTime: time.Now(),
		FileName:  fileName,
	}

	reader := &ProgressReader{
		Reader:     response.Body,
		Progress:   progress,
		UpdateChan: make(chan int64, 100),
	}

	go progress.UpdateProgress(reader.UpdateChan)

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	close(reader.UpdateChan)

	logx.Infof("文件下载完成: %s", filePath)
	return nil
}

type Progress struct {
	Total     int64
	Current   int64
	StartTime time.Time
	FileName  string
	LastBytes int64
	LastTime  time.Time
	Stopped   bool
}

type ProgressReader struct {
	Reader     io.Reader
	Progress   *Progress
	UpdateChan chan int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	if pr.Progress.Stopped {
		return 0, fmt.Errorf("下载已停止：速度过慢")
	}

	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.Progress.Current += int64(n)
		select {
		case pr.UpdateChan <- pr.Progress.Current:
		default:
		}
	}
	return n, err
}

func (p *Progress) UpdateProgress(updateChan chan int64) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	p.LastTime = time.Now()
	p.LastBytes = p.Current

	for {
		select {
		case current, ok := <-updateChan:
			if !ok {
				p.displayProgress(current)
				return
			}
			p.Current = current
		case <-ticker.C:
			now := time.Now()
			timeDiff := now.Sub(p.LastTime).Seconds()

			if timeDiff > 0 {
				speed := float64(p.Current-p.LastBytes) / timeDiff

				if speed < 1024 && now.Sub(p.StartTime).Hours() >= 1 {
					p.Stopped = true
					logx.Errorf("下载速度过慢，已停止下载: %s (速度: %.2f KB/s)", p.FileName, speed/1024)
					return
				}

				p.LastTime = now
				p.LastBytes = p.Current
			}

			p.displayProgress(p.Current)
		}
	}
}

func (p *Progress) displayProgress(current int64) {
	now := time.Now()
	duration := now.Sub(p.StartTime).Seconds()
	if duration > 0 {
		speed := float64(current) / duration / 1024 / 1024
		progress := float64(current) / float64(p.Total) * 100
		logx.Infof("下载进度 [%s]: %.2f%% (%.2f MB/s)", p.FileName, progress, speed)
	}
}

func (m *Auth115Manager) DownloadFile(fileID string) error {
	fileDetail, err := m.GetFileInfo(fileID)
	if err != nil {
		return fmt.Errorf("获取文件详情失败: %v", err)
	}

	downloadUrlResp, err := m.GetFileDownloadUrl(fileDetail.Data.PickCode)
	if err != nil {
		return fmt.Errorf("获取文件下载地址失败: %v", err)
	}

	var wg sync.WaitGroup

	for id, info := range downloadUrlResp.Data {
		logx.Infof("开始下载文件: %s, 大小: %d", info.FileName, info.FileSize)

		wg.Add(1)

		go func(id string, name string, size int64, url string) {
			defer wg.Done()

			if err := m.downloadFile(fileDetail.Data.FileName, id, name, size, url); err != nil {
				logx.Errorf("下载文件失败: %v", err)
			}
		}(id, info.FileName, info.FileSize, info.URL.URL)
	}

	wg.Wait()
	logx.Info("所有文件下载任务已完成")

	return nil
}

func (m *Auth115Manager) GetFolderFiles(folderID string, fileType int) (*types.GetFolderFilesResp, error) {
	accessToken, err := m.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}

	params := url.Values{}
	params.Set("cid", folderID)

	if fileType > 0 {
		params.Set("type", fmt.Sprintf("%d", fileType))
	}

	params.Set("limit", "1150")
	params.Set("offset", "0")
	params.Set("asc", "0")
	params.Set("o", "file_name")
	params.Set("stdir", "1")
	params.Set("cur", "1")

	apiURL := "https://proapi.115.com/open/ufile/files?" + params.Encode()
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result types.GetFolderFilesResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.State {
		return nil, fmt.Errorf("获取文件夹内文件列表失败: %s", result.Message)
	}

	return &result, nil
}

func (m *Auth115Manager) GetFileInfos() []FileInfo {
	m.fileMutex.Lock()
	defer m.fileMutex.Unlock()

	infos := make([]FileInfo, len(m.fileInfos))
	copy(infos, m.fileInfos)
	return infos
}

func (m *Auth115Manager) create115Folder(accessToken string) error {
	data := url.Values{}
	data.Set("pid", "0")
	data.Set("file_name", "115quick")

	apiURL := "https://proapi.115.com/open/folder/add"
	request, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var result struct {
		State   bool   `json:"state"`
		Message string `json:"message"`
		Code    int    `json:"code"`
		Data    struct {
			FileName string `json:"file_name"`
			FileID   string `json:"file_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.State {
		return fmt.Errorf("创建文件夹失败: %s", result.Message)
	}

	if result.Data.FileID != "" {
		m.svcCtx.SetFolderID(result.Data.FileID)
		logx.Info("成功创建115文件夹，文件夹ID: ", result.Data.FileID)
	} else {
		logx.Info("创建115文件夹成功但未返回文件夹ID，搜索文件夹ID")
	}
	return nil
}

func (m *Auth115Manager) search115Folder() {
	accessToken, err := m.GetAccessToken()
	if err != nil {
		logx.Errorf("获取访问令牌失败: %v", err)
		return
	}

	params := url.Values{}
	params.Set("search_value", "115quick")
	params.Set("limit", "20")
	params.Set("offset", "0")
	params.Set("fc", "1")

	apiURL := "https://proapi.115.com/open/ufile/search?" + params.Encode()
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logx.Errorf("创建请求失败: %v", err)
		return
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logx.Errorf("发送请求失败: %v", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logx.Errorf("读取响应失败: %v", err)
		return
	}

	var result types.Search115FolderResp
	if err := json.Unmarshal(body, &result); err != nil {
		logx.Errorf("解析响应失败: %v", err)
		return
	}

	if !result.State {
		logx.Errorf("搜索文件夹失败: %s", result.Message)
		return
	}

	if len(result.Data) > 0 {
		folderID := result.Data[0].FileID
		if folderID != "" {
			m.svcCtx.SetFolderID(folderID)
			logx.Infof("成功找到115文件夹，文件夹ID: %s", folderID)
		}
	} else {
		logx.Info("未找到115文件夹")
	}
}

func (m *Auth115Manager) InitFolder() {
	token, err := m.GetAccessToken()
	if err != nil {
		logx.Errorf("获取访问令牌失败，无法初始化文件夹: %v", err)
		return
	}
	if err := m.create115Folder(token); err != nil {
		logx.Errorf("创建115文件夹失败: %v", err)
		m.search115Folder()
	}
}
