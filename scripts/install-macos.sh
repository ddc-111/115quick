#!/bin/bash
set -e

# ============================================
# 115Quick macOS 一键安装脚本
# 直接执行即可安装并配置为系统服务
# ============================================

REPO="ddc-111/115quick"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/115quick"
DATA_DIR="$HOME/.local/share/115quick"
PLIST_NAME="com.115quick.server"
PLIST_PATH="$HOME/Library/LaunchAgents/${PLIST_NAME}.plist"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo "=========================================="
echo "  115Quick macOS 安装程序"
echo "=========================================="
echo ""

# 1. 检测架构
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64)  ARCH="arm64" ;;
    *)      echo "错误: 不支持的架构 $ARCH"; exit 1 ;;
esac
echo -e "系统架构: ${GREEN}darwin-${ARCH}${NC}"

# 2. 获取最新版本
echo "获取最新版本..."
VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
if [ -z "$VERSION" ]; then
    echo "错误: 无法获取版本信息"
    exit 1
fi
echo -e "最新版本: ${GREEN}v${VERSION}${NC}"

# 3. 下载
TMP_DIR=$(mktemp -d)
FILENAME="115quick-darwin-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${FILENAME}"
echo "下载中..."
curl -L -o "${TMP_DIR}/${FILENAME}" "$DOWNLOAD_URL"

# 4. 解压
tar -xzf "${TMP_DIR}/${FILENAME}" -C "$TMP_DIR"

# 5. 安装二进制文件
echo "安装到 ${INSTALL_DIR}..."
sudo mkdir -p "$INSTALL_DIR"
sudo cp "${TMP_DIR}/115quick" "${INSTALL_DIR}/115quick"
sudo chmod +x "${INSTALL_DIR}/115quick"

# 6. 创建目录
mkdir -p "$CONFIG_DIR"
mkdir -p "$DATA_DIR"

# 7. 创建默认配置文件
if [ ! -f "${CONFIG_DIR}/quick.yaml" ]; then
    cat > "${CONFIG_DIR}/quick.yaml" << EOF
Name: quick
Host: 0.0.0.0
Port: 8889
DBPath: ${DATA_DIR}/115quick.db
Auth115:
  DownloadPath: ${DATA_DIR}
  AccessToken: ""
  RefreshToken: ""
SMB:
  Enabled: false
  Host: ""
  Share: ""
  Username: ""
  Password: ""
  MountPoint: ""
EOF
    echo -e "配置文件已创建: ${GREEN}${CONFIG_DIR}/quick.yaml${NC}"
else
    echo "配置文件已存在，跳过"
fi

# 8. 创建 macOS 服务
cat > "$PLIST_PATH" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>${PLIST_NAME}</string>
    <key>ProgramArguments</key>
    <array>
        <string>${INSTALL_DIR}/115quick</string>
        <string>-f</string>
        <string>${CONFIG_DIR}/quick.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${DATA_DIR}/115quick.log</string>
    <key>StandardErrorPath</key>
    <string>${DATA_DIR}/115quick.error.log</string>
    <key>WorkingDirectory</key>
    <string>${DATA_DIR}</string>
</dict>
</plist>
EOF

# 9. 启动服务
launchctl unload "$PLIST_PATH" 2>/dev/null || true
launchctl load "$PLIST_PATH"
launchctl start "$PLIST_NAME"

# 10. 清理
rm -rf "$TMP_DIR"

# 11. 完成
echo ""
echo -e "${GREEN}=========================================="
echo "  安装完成!"
echo "==========================================${NC}"
echo ""
echo "服务已启动，开机自动运行"
echo ""
echo "配置文件: ${CONFIG_DIR}/quick.yaml"
echo "数据目录: ${DATA_DIR}"
echo "日志文件: ${DATA_DIR}/115quick.log"
echo ""
echo "访问地址: http://localhost:8889"
echo ""
echo "常用命令:"
echo "  115quick --version    查看版本"
echo "  115quick --help       查看帮助"
echo ""
echo "管理服务:"
echo "  launchctl stop com.115quick.server    停止"
echo "  launchctl start com.115quick.server   启动"
echo "  launchctl restart com.115quick.server 重启"
echo ""
