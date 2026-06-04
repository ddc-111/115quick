package Download

import (
	"context"
	"fmt"
	"path/filepath"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SMBDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSMBDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SMBDownloadLogic {
	return &SMBDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SMBDownloadLogic) SMBDownload(req *types.SMBDownloadReq) (resp *types.SMBDownloadResp, err error) {
	if !l.svcCtx.SMB.IsConnected() {
		return nil, fmt.Errorf("SMB server not connected, please configure and connect first")
	}

	if req.RemotePath == "" {
		return nil, fmt.Errorf("remote path is required")
	}

	// 确定本地保存路径
	localPath := req.LocalPath
	if localPath == "" {
		// 使用默认下载路径 + 远程文件名
		downloadPath := l.svcCtx.Auth115.GetDownloadPath()
		fileName := filepath.Base(req.RemotePath)
		localPath = filepath.Join(downloadPath, fileName)
	}

	// 下载文件
	if err := l.svcCtx.SMB.DownloadFile(req.RemotePath, localPath); err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	logx.Infof("SMB file downloaded: %s -> %s", req.RemotePath, localPath)

	return &types.SMBDownloadResp{
		Success:  true,
		Message:  "File downloaded successfully",
		FilePath: localPath,
	}, nil
}
