package Download

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFolderLogic {
	return &CreateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFolderLogic) CreateFolder(req *types.CreateFolderReq) (resp *types.CreateFolderResp, err error) {
	accessToken, err := l.svcCtx.Auth115.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}

	data := url.Values{}
	data.Set("pid", req.ParentID)
	data.Set("file_name", req.FolderName)

	apiURL := "https://proapi.115.com/open/folder/add"
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

	var result struct {
		State   bool   `json:"state"`
		Message string `json:"message"`
		Data    struct {
			FileID   string `json:"file_id"`
			FileName string `json:"file_name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.State {
		return nil, fmt.Errorf("创建文件夹失败: %s", result.Message)
	}

	resp = &types.CreateFolderResp{
		FileID:   result.Data.FileID,
		FileName: result.Data.FileName,
	}

	return resp, nil
}
