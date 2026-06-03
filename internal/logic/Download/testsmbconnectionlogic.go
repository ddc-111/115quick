package Download

import (
	"context"
	"fmt"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"
	"115Quick_server/internal/utils/smb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TestSMBConnectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestSMBConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestSMBConnectionLogic {
	return &TestSMBConnectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestSMBConnectionLogic) TestSMBConnection(req *types.TestSMBConnectionReq) (resp *types.NilResp, err error) {
	if req.Host == "" || req.Share == "" {
		return nil, fmt.Errorf("SMB host and share are required")
	}

 smbCfg := &smb.SMBConfig{
		Host:     req.Host,
		Share:    req.Share,
		Username: req.Username,
		Password: req.Password,
	}

	manager := smb.NewSMBManager()
	if err := manager.TestConnection(smbCfg); err != nil {
		return nil, fmt.Errorf("SMB connection test failed: %v", err)
	}

	return &types.NilResp{Success: true, Message: "SMB connection test successful"}, nil
}
