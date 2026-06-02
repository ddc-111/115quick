package Download

import (
	"context"
	"time"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取云下载任务列表
func NewGetCloudTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudTasksLogic {
	return &GetCloudTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCloudTasksLogic) GetCloudTasks(req *types.CloudTasksReq) (resp *types.CloudTasksResp, err error) {
	tasks := l.svcCtx.Auth115.GetCloudTasks()

	items := make([]types.CloudTaskItem, len(tasks))
	for i, task := range tasks {
		addTime := ""
		if task.AddTime > 0 {
			addTime = time.Unix(task.AddTime, 0).Format("2006-01-02 15:04:05")
		}

		items[i] = types.CloudTaskItem{
			InfoHash:     task.InfoHash,
			Name:         task.Name,
			Size:         task.Size,
			Status:       task.Status,
			PercentDone:  task.PercentDone,
			RateDownload: task.RateDownload,
			AddTime:      addTime,
			URL:          task.URL,
		}
	}

	return &types.CloudTasksResp{
		Tasks: items,
		Count: len(items),
	}, nil
}
