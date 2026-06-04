package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveDownloadTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除待下载任务
func NewRemoveDownloadTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveDownloadTaskLogic {
	return &RemoveDownloadTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveDownloadTaskLogic) RemoveDownloadTask(req *types.RemoveDownloadTaskReq) (resp *types.NilResp, err error) {
	if req.URL == "" {
		return nil, fmt.Errorf("URL不能为空")
	}

	if err := l.svcCtx.Auth115.RemovePendingTask(req.URL); err != nil {
		return nil, fmt.Errorf("删除任务失败: %v", err)
	}

	return &types.NilResp{}, nil
}
