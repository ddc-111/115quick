package Download

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServerLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerLogsLogic {
	return &GetServerLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServerLogsLogic) GetServerLogs(req *types.GetServerLogsReq) (resp *types.GetServerLogsResp, err error) {
	resp = &types.GetServerLogsResp{
		Logs: []string{},
	}

	logType := req.Type
	if logType == "" {
		logType = "stdout"
	}

	lines := req.Lines
	if lines <= 0 {
		lines = 200
	}

	// 获取日志文件路径
	logDir := filepath.Dir(l.svcCtx.Config.DBPath)
	var logFile string
	if logType == "stderr" {
		logFile = filepath.Join(logDir, "115quick.error.log")
	} else {
		logFile = filepath.Join(logDir, "115quick.log")
	}

	// 读取日志文件
	content, err := os.ReadFile(logFile)
	if err != nil {
		resp.Logs = []string{"无法读取日志文件: " + logFile + " - " + err.Error()}
		return resp, nil
	}

	// 获取最后N行
	logLines := strings.Split(string(content), "\n")
	if len(logLines) > lines {
		logLines = logLines[len(logLines)-lines:]
	}

	// 过滤空行
	for _, line := range logLines {
		if strings.TrimSpace(line) != "" {
			resp.Logs = append(resp.Logs, line)
		}
	}

	return resp, nil
}
