package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"115Quick_server/internal/config"
	"115Quick_server/internal/handler"
	"115Quick_server/internal/svc"
	"115Quick_server/internal/version"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var (
	configFile = flag.String("f", "etc/quick.yaml", "the config file")
	showVer    = flag.Bool("version", false, "show version")
)

func main() {
	flag.Parse()

	if *showVer {
		fmt.Printf("115Quick v%s (commit: %s, built: %s)\n", version.Version, version.GitCommit, version.BuildTime)
		os.Exit(0)
	}

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	ctx.Auth115.InitFolder()
	ctx.StarDownloadLinkWorker()
	handler.RegisterHandlers(server, ctx)

	// 添加日志查看端点
	http.HandleFunc("/v1/Download/api/getServerLogs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		logType := r.URL.Query().Get("type")
		if logType == "" {
			logType = "stdout"
		}

		linesStr := r.URL.Query().Get("lines")
		lines := 200
		if l, err := strconv.Atoi(linesStr); err == nil && l > 0 {
			lines = l
		}

		// 获取日志文件路径
		logDir := filepath.Dir(c.DBPath)
		var logFile string
		if logType == "stderr" {
			logFile = filepath.Join(logDir, "115quick.error.log")
		} else {
			logFile = filepath.Join(logDir, "115quick.log")
		}

		// 读取日志文件
		content, err := os.ReadFile(logFile)
		if err != nil {
			fmt.Fprintf(w, `{"code":200,"data":{"logs":["无法读取日志文件: %s"],"error":"%s"}}`, logFile, err.Error())
			return
		}

		// 获取最后N行
		logLines := strings.Split(string(content), "\n")
		if len(logLines) > lines {
			logLines = logLines[len(logLines)-lines:]
		}

		// 构建JSON响应
		fmt.Fprintf(w, `{"code":200,"data":{"logs":[`)
		for i, line := range logLines {
			if i > 0 {
				fmt.Fprintf(w, `,`)
			}
			// 转义JSON特殊字符
			line = strings.ReplaceAll(line, `\`, `\\`)
			line = strings.ReplaceAll(line, `"`, `\"`)
			line = strings.ReplaceAll(line, "\n", `\n`)
			line = strings.ReplaceAll(line, "\r", `\r`)
			line = strings.ReplaceAll(line, "\t", `\t`)
			fmt.Fprintf(w, `"%s"`, line)
		}
		fmt.Fprintf(w, `]}}`)
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
