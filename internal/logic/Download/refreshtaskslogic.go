package Download

import (
	"context"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手动刷新任务列表
func NewRefreshTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTasksLogic {
	return &RefreshTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTasksLogic) RefreshTasks(req *types.RefreshTasksReq) (resp *types.NilResp, err error) {
	l.svcCtx.Auth115.RefreshCloudTasks()
	return &types.NilResp{}, nil
}
