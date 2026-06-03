#!/bin/bash
set -e

# ============================================
# 115Quick 版本检测与更新脚本
# ============================================

REPO="ddc-111/115quick"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/115quick"
DATA_DIR="$HOME/.local/share/115quick"
PLIST_NAME="com.115quick.server"
PLIST_PATH="$HOME/Library/LaunchAgents/${PLIST_NAME}.plist"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 获取当前安装版本
get_installed_version() {
    if [ -f "${INSTALL_DIR}/115quick" ]; then
        "${INSTALL_DIR}/115quick" --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "unknown"
    else
        echo "none"
    fi
}

# 获取最新版本
get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/'
}

# 检测架构
detect_arch() {
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64) echo "amd64" ;;
        arm64)  echo "arm64" ;;
        *)      echo "unknown" ;;
    esac
}

# 比较版本
version_gt() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ] && [ "$1" != "$2" ]
}

# 显示状态
show_status() {
    echo ""
    echo "=========================================="
    echo "  115Quick 状态"
    echo "=========================================="
    echo ""

    local installed
    installed=$(get_installed_version)
    echo "安装版本: $installed"
    echo "安装路径: ${INSTALL_DIR}/115quick"
    echo "配置目录: ${CONFIG_DIR}"
    echo "数据目录: ${DATA_DIR}"

    if [ -f "$PLIST_PATH" ]; then
        if launchctl list | grep -q "$PLIST_NAME"; then
            echo -e "服务状态: ${GREEN}运行中${NC}"
        else
            echo -e "服务状态: ${YELLOW}已安装未运行${NC}"
        fi
    else
        echo -e "服务状态: ${RED}未安装${NC}"
    fi

    if lsof -i :8889 &>/dev/null; then
        echo -e "端口 8889: ${GREEN}已占用${NC}"
    else
        echo -e "端口 8889: ${YELLOW}未占用${NC}"
    fi
}

# 检查更新
check_update() {
    echo ""
    echo "检查更新..."
    echo ""

    local installed
    installed=$(get_installed_version)
    local latest
    latest=$(get_latest_version)

    if [ -z "$latest" ]; then
        echo -e "${RED}无法获取版本信息${NC}"
        exit 1
    fi

    echo "当前版本: $installed"
    echo "最新版本: $latest"

    if [ "$installed" = "none" ]; then
        echo -e "${YELLOW}未安装 115Quick${NC}"
        echo ""
        read -p "是否安装？(Y/n): " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            do_install "$latest"
        fi
    elif version_gt "$latest" "$installed"; then
        echo ""
        echo -e "${YELLOW}发现新版本 v${latest}${NC}"
        echo ""
        read -p "是否更新？(Y/n): " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            do_update "$latest"
        fi
    else
        echo ""
        echo -e "${GREEN}已是最新版本${NC}"
    fi
}

# 执行安装
do_install() {
    local version="$1"
    local arch
    arch=$(detect_arch)
    local filename="115quick-darwin-${arch}.tar.gz"
    local tmp_dir
    tmp_dir=$(mktemp -d)

    echo "下载 v${version}..."
    curl -L -o "${tmp_dir}/${filename}" "https://github.com/${REPO}/releases/download/v${version}/${filename}"
    tar -xzf "${tmp_dir}/${filename}" -C "$tmp_dir"

    sudo mkdir -p "$INSTALL_DIR"
    sudo cp "${tmp_dir}/115quick" "${INSTALL_DIR}/115quick"
    sudo chmod +x "${INSTALL_DIR}/115quick"

    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"

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
    fi

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

    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    launchctl load "$PLIST_PATH"
    launchctl start "$PLIST_NAME"

    rm -rf "$tmp_dir"

    echo ""
    echo -e "${GREEN}安装完成!${NC}"
}

# 执行更新
do_update() {
    local version="$1"
    local arch
    arch=$(detect_arch)
    local filename="115quick-darwin-${arch}.tar.gz"
    local tmp_dir
    tmp_dir=$(mktemp -d)

    echo "停止服务..."
    launchctl stop "$PLIST_NAME" 2>/dev/null || true

    echo "下载 v${version}..."
    curl -L -o "${tmp_dir}/${filename}" "https://github.com/${REPO}/releases/download/v${version}/${filename}"
    tar -xzf "${tmp_dir}/${filename}" -C "$tmp_dir"

    echo "更新二进制文件..."
    sudo cp "${tmp_dir}/115quick" "${INSTALL_DIR}/115quick"
    sudo chmod +x "${INSTALL_DIR}/115quick"

    echo "启动服务..."
    launchctl start "$PLIST_NAME"

    rm -rf "$tmp_dir"

    echo ""
    echo -e "${GREEN}更新完成! 当前版本: v${version}${NC}"
}

# 主逻辑
case "${1:-}" in
    --status|-s)
        show_status
        ;;
    --help|-h)
        echo ""
        echo "115Quick 版本检测与更新工具"
        echo ""
        echo "用法:"
        echo "  115quick-update          检查更新"
        echo "  115quick-update --status 显示状态"
        echo "  115quick-update --help   显示帮助"
        echo ""
        ;;
    *)
        check_update
        ;;
esac
