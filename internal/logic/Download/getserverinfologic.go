package Download

import (
	"context"
	"sort"
	"time"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"
	"115Quick_server/internal/version"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerInfoLogic {
	return &GetServerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerInfoLogic) GetServerInfo(req *types.ServerInfoReq) (resp *types.ServerInfoResp, err error) {
	fileInfos := l.svcCtx.Auth115.GetFileInfos()
	resp = &types.ServerInfoResp{
		DownFileInfoList: make([]types.DownFileInfo, len(fileInfos)),
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].AddTime.After(fileInfos[j].AddTime)
	})
	for i, info := range fileInfos {
		resp.DownFileInfoList[i] = types.DownFileInfo{
			Url:     info.URL,
			AddTime: info.AddTime.Format("2006-01-02 15:04:05"),
		}
	}
	if len(fileInfos) == 0 {
		var downFileInfoList []types.DownFileInfo
		downFileInfoList = append(downFileInfoList, types.DownFileInfo{Url: "暂无等待下载链接 服务器正常运行中", AddTime: time.Now().Format("2006-01-02 15:04:05")})
		resp.DownFileInfoList = downFileInfoList
	}

	resp.AuthData = types.AuthData{
		AuthStatus: "正常",
		ExpireAt:   "",
	}
	resp.Mode = l.svcCtx.Mode

	versionInfo := version.GetVersionInfo()
	resp.Version = types.VersionInfo{
		Version:       versionInfo.Version,
		GitCommit:     versionInfo.GitCommit,
		BuildTime:     versionInfo.BuildTime,
		LatestVersion: versionInfo.LatestVersion,
		UpdateURL:     versionInfo.UpdateURL,
		HasUpdate:     versionInfo.HasUpdate,
	}

	return
}
