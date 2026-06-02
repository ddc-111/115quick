package request115

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// API接口常量
const (
	API_BASE_URL        = "https://proapi.115.com" // 115云服务API基础URL
	API_GET_FILE_LIST   = "/open/ufile/files"
	API_GET_FOLDER_INFO = "/open/folder/get_info"
	API_RENAME_FILE     = "/open/ufile/update"
	API_DELETE_FILE     = "/open/ufile/delete" // 假设删除文件的API路径
	API_CREATE_FOLDER   = "/open/folder/add"
)

var (
	AccessToken = "" // 全局AccessToken，需要在应用启动时设置
)

// 设置AccessToken
func SetAccessToken(token string) {
	AccessToken = token
}

// 文件列表返回结构
type FileListResponse struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    []File `json:"data"`
	Count   int    `json:"count"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

// 文件信息结构
type File struct {
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size,string"` // 使用int64类型并添加string标签，以便于JSON解析字符串格式的数字
	FileCategory string `json:"file_category"`    // "0"表示文件夹，"1"表示文件
	PickCode     string `json:"pick_code"`
	ParentID     string `json:"parent_id"`
	Sha1         string `json:"sha1"`
	Ptime        string `json:"ptime"` // 上传时间
	Utime        string `json:"utime"` // 修改时间
}

// 文件夹信息返回结构
type FolderInfoResponse struct {
	State   bool       `json:"state"`
	Message string     `json:"message"`
	Code    int        `json:"code"`
	Data    FolderInfo `json:"data"`
}

// 文件夹信息结构
type FolderInfo struct {
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	FileCategory string `json:"file_category"` // "0"表示文件夹
	Count        string `json:"count"`         // 包含文件总数量
	Size         string `json:"size"`          // 文件夹总大小
	FolderCount  string `json:"folder_count"`  // 包含文件夹总数量
}

// 基础响应结构
type BaseResponse struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// 重命名响应结构
type RenameResponse struct {
	BaseResponse
	Data struct {
		FileName string `json:"file_name"`
		Star     string `json:"star"`
	} `json:"data"`
}

// 获取文件列表
func GetFileList(folderID string, showDir int, limit int) (*FileListResponse, error) {
	apiURL := API_BASE_URL + API_GET_FILE_LIST

	// 构建请求参数
	params := url.Values{}
	params.Add("cid", folderID)
	params.Add("show_dir", strconv.Itoa(showDir))
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", "0") // 从0开始

	// 发送GET请求
	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// 添加授权头
	req.Header.Add("Authorization", "Bearer "+AccessToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var fileListResp FileListResponse
	err = json.Unmarshal(body, &fileListResp)
	if err != nil {
		return nil, err
	}

	if !fileListResp.State {
		return nil, fmt.Errorf("API错误: %s (代码: %d)", fileListResp.Message, fileListResp.Code)
	}

	return &fileListResp, nil
}

// 获取文件夹信息
func GetFolderInfo(folderID string) (*FolderInfo, error) {
	apiURL := API_BASE_URL + API_GET_FOLDER_INFO

	// 构建请求参数
	params := url.Values{}
	params.Add("file_id", folderID)

	// 发送GET请求
	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// 添加授权头
	req.Header.Add("Authorization", "Bearer "+AccessToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var folderInfoResp FolderInfoResponse
	err = json.Unmarshal(body, &folderInfoResp)
	if err != nil {
		return nil, err
	}

	if !folderInfoResp.State {
		return nil, fmt.Errorf("API错误: %s (代码: %d)", folderInfoResp.Message, folderInfoResp.Code)
	}

	return &folderInfoResp.Data, nil
}

// 重命名文件或文件夹
func RenameFile(fileID string, newName string) error {
	apiURL := API_BASE_URL + API_RENAME_FILE

	// 构建表单数据
	formData := url.Values{}
	formData.Add("file_id", fileID)
	formData.Add("file_name", newName)

	// 发送POST请求
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	// 添加头信息
	req.Header.Add("Authorization", "Bearer "+AccessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析响应
	var renameResp RenameResponse
	err = json.Unmarshal(body, &renameResp)
	if err != nil {
		return err
	}

	if !renameResp.State {
		return fmt.Errorf("API错误: %s (代码: %d)", renameResp.Message, renameResp.Code)
	}

	return nil
}

// 删除文件
func DeleteFile(fileID string) error {
	apiURL := API_BASE_URL + API_DELETE_FILE

	// 构建表单数据
	formData := url.Values{}
	formData.Add("file_id", fileID)

	// 发送POST请求
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	// 添加头信息
	req.Header.Add("Authorization", "Bearer "+AccessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析响应
	var baseResp BaseResponse
	err = json.Unmarshal(body, &baseResp)
	if err != nil {
		return err
	}

	if !baseResp.State {
		return fmt.Errorf("API错误: %s (代码: %d)", baseResp.Message, baseResp.Code)
	}

	return nil
}
