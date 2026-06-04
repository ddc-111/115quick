package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearCompletedTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearCompletedTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCompletedTasksLogic {
	return &ClearCompletedTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearCompletedTasksLogic) ClearCompletedTasks() (resp *types.NilResp, err error) {
	// 清空已完成的任务历史
	if err := l.svcCtx.Store.ClearCompletedTaskHistory(); err != nil {
		return nil, fmt.Errorf("清空已完成任务失败: %v", err)
	}

	// 清空已完成的文件信息
	if err := l.svcCtx.Store.ClearCompletedTasks(); err != nil {
		return nil, fmt.Errorf("清空已完成任务失败: %v", err)
	}

	return &types.NilResp{Success: true, Message: "已完成任务已清空"}, nil
}
