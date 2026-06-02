package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置115 Token
func NewSetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTokenLogic {
	return &SetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetTokenLogic) SetToken(req *types.SetTokenReq) (resp *types.NilResp, err error) {
	if req.AccessToken == "" || req.RefreshToken == "" {
		return nil, fmt.Errorf("AccessToken和RefreshToken不能为空")
	}

	if err := l.svcCtx.Auth115.SetToken(req.AccessToken, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("设置Token失败: %v", err)
	}

	return &types.NilResp{}, nil
}
