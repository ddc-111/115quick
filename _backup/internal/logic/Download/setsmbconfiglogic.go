package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"
	"115Quick_server/internal/utils/smb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetSMBConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetSMBConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSMBConfigLogic {
	return &SetSMBConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetSMBConfigLogic) SetSMBConfig(req *types.SetSMBConfigReq) (resp *types.NilResp, err error) {
	if req.Enabled {
		if req.Host == "" || req.Share == "" {
			return nil, fmt.Errorf("SMB host and share are required when enabling SMB")
		}

		smbCfg := &smb.SMBConfig{
			Host:     req.Host,
			Share:    req.Share,
			Username: req.Username,
			Password: req.Password,
		}

		if err := l.svcCtx.SMB.Connect(smbCfg); err != nil {
			return nil, fmt.Errorf("failed to connect to SMB server: %v", err)
		}

		logx.Infof("SMB server connected: %s/%s", req.Host, req.Share)
	} else {
		if l.svcCtx.SMB.IsConnected() {
			if err := l.svcCtx.SMB.Disconnect(); err != nil {
				return nil, fmt.Errorf("failed to disconnect SMB: %v", err)
			}
		}

		logx.Info("SMB disabled, disconnected")
	}

	cfg := svc.SMBConfigData{
		Enabled:  req.Enabled,
		Host:     req.Host,
		Share:    req.Share,
		Username: req.Username,
		Password: req.Password,
	}
	if err := l.svcCtx.Store.SaveSMBConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to save SMB config: %v", err)
	}

	return &types.NilResp{Success: true, Message: "SMB configuration updated"}, nil
}
