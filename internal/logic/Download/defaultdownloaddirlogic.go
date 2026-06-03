package Download

import (
	"context"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultDownloadDirLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultDownloadDirLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultDownloadDirLogic {
	return &SetDefaultDownloadDirLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultDownloadDirLogic) SetDefaultDownloadDir(req *types.SetDefaultDownloadDirReq) (resp *types.NilResp, err error) {
	err = l.svcCtx.Store.SetConfig("default_download_folder_id", req.FolderID)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.Store.SetConfig("default_download_folder_name", req.FolderName)
	if err != nil {
		return nil, err
	}

	l.svcCtx.SetFolderID(req.FolderID)

	resp = &types.NilResp{
		Success: true,
		Message: "默认下载目录已设置",
	}

	return resp, nil
}

type GetDefaultDownloadDirLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDefaultDownloadDirLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDefaultDownloadDirLogic {
	return &GetDefaultDownloadDirLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDefaultDownloadDirLogic) GetDefaultDownloadDir(req *types.GetDefaultDownloadDirReq) (resp *types.GetDefaultDownloadDirResp, err error) {
	folderID, _ := l.svcCtx.Store.GetConfig("default_download_folder_id")
	folderName, _ := l.svcCtx.Store.GetConfig("default_download_folder_name")

	resp = &types.GetDefaultDownloadDirResp{
		FolderID:   folderID,
		FolderName: folderName,
	}

	return resp, nil
}
