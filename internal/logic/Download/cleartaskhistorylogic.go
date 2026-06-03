package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearTaskHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearTaskHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearTaskHistoryLogic {
	return &ClearTaskHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearTaskHistoryLogic) ClearTaskHistory() (resp *types.NilResp, err error) {
	if err := l.svcCtx.Store.ClearTaskHistory(); err != nil {
		return nil, fmt.Errorf("清空任务历史失败: %v", err)
	}

	return &types.NilResp{Success: true, Message: "任务历史已清空"}, nil
}
