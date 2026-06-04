package Download

import (
	"context"
	"fmt"
	"time"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SMBBrowseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSMBBrowseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SMBBrowseLogic {
	return &SMBBrowseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SMBBrowseLogic) SMBBrowse(req *types.SMBBrowseReq) (resp *types.SMBBrowseResp, err error) {
	if !l.svcCtx.SMB.IsConnected() {
		return nil, fmt.Errorf("SMB server not connected, please configure and connect first")
	}

	path := req.Path
	if path == "" {
		path = "."
	}

	files, err := l.svcCtx.SMB.ListFiles(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list SMB files: %v", err)
	}

	fileItems := make([]types.SMBFileItem, 0, len(files))
	for _, f := range files {
		fileItems = append(fileItems, types.SMBFileItem{
			Name:    f.Name,
			Size:    f.Size,
			IsDir:   f.IsDir,
			ModTime: f.ModTime.Format(time.DateTime),
		})
	}

	return &types.SMBBrowseResp{
		Files:       fileItems,
		CurrentPath: path,
	}, nil
}
