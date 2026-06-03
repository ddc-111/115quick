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
		if req.MountPoint == "" {
			return nil, fmt.Errorf("mount point is required when enabling SMB")
		}

		smbCfg := &smb.SMBConfig{
			Host:       req.Host,
			Share:      req.Share,
			Username:   req.Username,
			Password:   req.Password,
			MountPoint: req.MountPoint,
		}

		if err := l.svcCtx.SMB.Mount(smbCfg); err != nil {
			return nil, fmt.Errorf("failed to mount SMB share: %v", err)
		}

		l.svcCtx.Auth115.SetDownloadPath(req.MountPoint)
		logx.Infof("SMB share mounted at %s", req.MountPoint)
	} else {
		if l.svcCtx.SMB.IsMounted() {
			if err := l.svcCtx.SMB.Unmount(); err != nil {
				return nil, fmt.Errorf("failed to unmount SMB share: %v", err)
			}
		}

		l.svcCtx.Auth115.SetDownloadPath(l.svcCtx.Config.Auth115.DownloadPath)
		logx.Info("SMB disabled, using default download path")
	}

	cfg := svc.SMBConfigData{
		Enabled:    req.Enabled,
		Host:       req.Host,
		Share:      req.Share,
		Username:   req.Username,
		Password:   req.Password,
		MountPoint: req.MountPoint,
	}
	if err := l.svcCtx.Store.SaveSMBConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to save SMB config: %v", err)
	}

	return &types.NilResp{Success: true, Message: "SMB configuration updated"}, nil
}
