package Download

import (
	"context"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDownloadModeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置下载模式
func NewSetDownloadModeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDownloadModeLogic {
	return &SetDownloadModeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDownloadModeLogic) SetDownloadMode(req *types.SetDownloadModeReq) (resp *types.NilResp, err error) {
	if req.Mode == 1 {
		l.svcCtx.Mode = 1
	} else {
		l.svcCtx.Mode = 0
	}
	return
}
