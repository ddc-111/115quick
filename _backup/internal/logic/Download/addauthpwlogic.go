package Download

import (
	"context"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAuthPwLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAuthPwLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAuthPwLogic {
	return &AddAuthPwLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAuthPwLogic) AddAuthPw(req *types.AddAuthPwReq) (resp *types.AddAuthPwResp, err error) {
	return
}
