package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任务历史
func NewGetTaskHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskHistoryLogic {
	return &GetTaskHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskHistoryLogic) GetTaskHistory(req *types.TaskHistoryReq) (resp *types.TaskHistoryResp, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	items, total, err := l.svcCtx.Store.GetTaskHistory(req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取任务历史失败: %v", err)
	}

	historyItems := make([]types.TaskHistoryItem, len(items))
	for i, item := range items {
		completeTime := ""
		if item.CompleteTime != nil {
			completeTime = item.CompleteTime.Format("2006-01-02 15:04:05")
		}

		historyItems[i] = types.TaskHistoryItem{
			ID:           item.ID,
			URL:          item.URL,
			Name:         item.Name,
			Size:         item.Size,
			Status:       item.Status,
			Progress:     item.Progress,
			ErrorMsg:     item.ErrorMsg,
			AddTime:      item.AddTime.Format("2006-01-02 15:04:05"),
			CompleteTime: completeTime,
		}
	}

	return &types.TaskHistoryResp{
		Items:    historyItems,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
