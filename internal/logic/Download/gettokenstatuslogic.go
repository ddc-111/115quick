package Download

import (
	"context"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取Token状态
func NewGetTokenStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenStatusLogic {
	return &GetTokenStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenStatusLogic) GetTokenStatus(req *types.TokenStatusReq) (resp *types.TokenStatusResp, err error) {
	configured, valid, expiresAt, message := l.svcCtx.Auth115.GetTokenStatus()

	return &types.TokenStatusResp{
		Configured: configured,
		Valid:      valid,
		ExpiresAt:  expiresAt,
		Message:    message,
	}, nil
}
