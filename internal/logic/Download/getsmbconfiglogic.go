package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSMBConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSMBConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSMBConfigLogic {
	return &GetSMBConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSMBConfigLogic) GetSMBConfig(req *types.GetSMBConfigReq) (resp *types.GetSMBConfigResp, err error) {
	cfg, err := l.svcCtx.Store.GetSMBConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get SMB config: %v", err)
	}

	return &types.GetSMBConfigResp{
		Enabled:    cfg.Enabled,
		Host:       cfg.Host,
		Share:      cfg.Share,
		Username:   cfg.Username,
		Password:   cfg.Password,
		MountPoint: cfg.MountPoint,
		IsMounted:  l.svcCtx.SMB.IsMounted(),
	}, nil
}
