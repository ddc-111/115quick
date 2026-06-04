package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCloudFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudFilesLogic {
	return &GetCloudFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCloudFilesLogic) GetCloudFiles(req *types.GetCloudFilesReq) (resp *types.GetCloudFilesResp, err error) {
	folderID := req.FolderID
	if folderID == "" {
		folderID = "0"
	}

	result, err := l.svcCtx.Auth115.GetFolderFiles(folderID, req.FileType)
	if err != nil {
		return nil, fmt.Errorf("获取云文件列表失败: %v", err)
	}

	files := make([]types.CloudFile, 0)
	for _, item := range result.Data {
		isDir := item.FC == "0"
		files = append(files, types.CloudFile{
			FileID:     item.FID,
			FileName:   item.FN,
			FileSize:   item.FS,
			IsDir:      isDir,
			ParentID:   item.PID,
			UpdateTime: fmt.Sprintf("%d", item.UPT),
			PickCode:   item.PC,
		})
	}

	resp = &types.GetCloudFilesResp{
		Files:     files,
		CurrentID: folderID,
		ParentID:  folderID,
	}

	return resp, nil
}
