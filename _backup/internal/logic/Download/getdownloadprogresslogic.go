package Download

import (
	"context"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDownloadProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取下载进度
func NewGetDownloadProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDownloadProgressLogic {
	return &GetDownloadProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDownloadProgressLogic) GetDownloadProgress(req *types.DownloadProgressReq) (resp *types.DownloadProgressResp, err error) {
	progresses := l.svcCtx.Auth115.GetDownloadProgress()

	items := make([]types.DownloadProgressItem, len(progresses))
	for i, p := range progresses {
		items[i] = types.DownloadProgressItem{
			FileName:   p.FileName,
			FileSize:   p.FileSize,
			Downloaded: p.Downloaded,
			Speed:      p.Speed,
			Percent:    p.Percent,
			StartTime:  p.StartTime.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.DownloadProgressResp{
		Downloads: items,
	}, nil
}
