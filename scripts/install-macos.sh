#!/bin/bash

set -e

# ============================================
# 115Quick macOS 安装/更新脚本
# ============================================

REPO="AnJiaHa/115Quick_server"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/115quick"
DATA_DIR="$HOME/.local/share/115quick"
PLIST_NAME="com.115quick.server"
PLIST_PATH="$HOME/Library/LaunchAgents/${PLIST_NAME}.plist"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 检测系统架构
detect_arch() {
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64)
            echo "amd64"
            ;;
        arm64)
            echo "arm64"
            ;;
        *)
            error "不支持的架构: $arch"
            ;;
    esac
}

# 获取最新版本
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "无法获取最新版本信息"
    fi
    echo "$version"
}

# 获取当前安装版本
get_installed_version() {
    if [ -f "${INSTALL_DIR}/115quick" ]; then
        "${INSTALL_DIR}/115quick" --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' || echo "unknown"
    else
        echo "none"
    fi
}

# 下载文件
download() {
    local url="$1"
    local output="$2"
    info "下载: $url"
    if command -v curl &> /dev/null; then
        curl -L -o "$output" "$url"
    elif command -v wget &> /dev/null; then
        wget -O "$output" "$url"
    else
        error "需要 curl 或 wget"
    fi
}

# 安装二进制文件
install_binary() {
    local version="$1"
    local arch
    arch=$(detect_arch)
    local filename="115quick-darwin-${arch}.tar.gz"
    local download_url="https://github.com/${REPO}/releases/download/v${version}/${filename}"
    local tmp_dir
    tmp_dir=$(mktemp -d)

    info "正在下载 115Quick v${version} (darwin-${arch})..."
    download "$download_url" "${tmp_dir}/${filename}"

    info "解压文件..."
    tar -xzf "${tmp_dir}/${filename}" -C "$tmp_dir"

    info "安装到 ${INSTALL_DIR}..."
    sudo mkdir -p "$INSTALL_DIR"
    sudo cp "${tmp_dir}/115quick" "${INSTALL_DIR}/115quick"
    sudo chmod +x "${INSTALL_DIR}/115quick"

    # 创建配置目录
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"

    # 如果配置文件不存在，复制默认配置
    if [ ! -f "${CONFIG_DIR}/quick.yaml" ]; then
        if [ -f "${tmp_dir}/quick.yaml" ]; then
            cp "${tmp_dir}/quick.yaml" "${CONFIG_DIR}/quick.yaml"
            # 更新配置中的路径
            sed -i '' "s|DBPath: data/115quick.db|DBPath: ${DATA_DIR}/115quick.db|g" "${CONFIG_DIR}/quick.yaml"
            sed -i '' "s|DownloadPath: data|DownloadPath: ${DATA_DIR}|g" "${CONFIG_DIR}/quick.yaml"
            success "配置文件已创建: ${CONFIG_DIR}/quick.yaml"
        fi
    else
        info "配置文件已存在，跳过"
    fi

    # 清理
    rm -rf "$tmp_dir"

    success "115Quick v${version} 安装完成"
}

# 创建 launchd 服务
create_service() {
    info "创建 macOS 服务..."

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

    success "服务配置文件已创建: ${PLIST_PATH}"
}

# 启动服务
start_service() {
    info "启动服务..."
    launchctl load "$PLIST_PATH" 2>/dev/null || true
    launchctl start "$PLIST_NAME" 2>/dev/null || true
    success "服务已启动"
}

# 停止服务
stop_service() {
    info "停止服务..."
    launchctl stop "$PLIST_NAME" 2>/dev/null || true
    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    success "服务已停止"
}

# 卸载服务
remove_service() {
    stop_service
    rm -f "$PLIST_PATH"
    success "服务已移除"
}

# 卸载程序
uninstall() {
    info "卸载 115Quick..."

    # 停止并移除服务
    if [ -f "$PLIST_PATH" ]; then
        remove_service
    fi

    # 删除二进制文件
    if [ -f "${INSTALL_DIR}/115quick" ]; then
        sudo rm -f "${INSTALL_DIR}/115quick"
        success "已删除 ${INSTALL_DIR}/115quick"
    fi

    # 询问是否删除配置和数据
    echo ""
    read -p "是否删除配置文件和数据？(y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        rm -rf "$DATA_DIR"
        success "已删除配置和数据"
    else
        info "保留配置和数据: ${CONFIG_DIR}, ${DATA_DIR}"
    fi

    success "卸载完成"
}

# 检查更新
check_update() {
    local installed_version
    installed_version=$(get_installed_version)
    local latest_version
    latest_version=$(get_latest_version)

    echo ""
    echo "当前版本: ${installed_version}"
    echo "最新版本: ${latest_version}"

    if [ "$installed_version" = "$latest_version" ]; then
        success "已是最新版本"
    else
        warn "有新版本可用"
        echo ""
        read -p "是否更新到 v${latest_version}？(Y/n): " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            local was_running=false
            if [ -f "$PLIST_PATH" ]; then
                was_running=true
                stop_service
            fi
            install_binary "$latest_version"
            if [ "$was_running" = true ]; then
                start_service
            fi
        fi
    fi
}

# 显示状态
show_status() {
    echo ""
    echo "=== 115Quick 状态 ==="
    echo ""

    # 版本信息
    local installed_version
    installed_version=$(get_installed_version)
    echo "安装版本: ${installed_version}"
    echo "安装路径: ${INSTALL_DIR}/115quick"
    echo "配置目录: ${CONFIG_DIR}"
    echo "数据目录: ${DATA_DIR}"
    echo ""

    # 服务状态
    if [ -f "$PLIST_PATH" ]; then
        if launchctl list | grep -q "$PLIST_NAME"; then
            echo "服务状态: ${GREEN}运行中${NC}"
        else
            echo "服务状态: ${YELLOW}已安装未运行${NC}"
        fi
    else
        echo "服务状态: ${RED}未安装${NC}"
    fi

    # 端口检查
    if lsof -i :8889 &>/dev/null; then
        echo "端口 8889: ${GREEN}已占用${NC}"
    else
        echo "端口 8889: ${YELLOW}未占用${NC}"
    fi
    echo ""
}

# 显示帮助
show_help() {
    echo ""
    echo "115Quick macOS 安装管理工具"
    echo ""
    echo "用法: $0 <命令>"
    echo ""
    echo "命令:"
    echo "  install    安装 115Quick（如果已安装则更新）"
    echo "  update     检查并更新到最新版本"
    echo "  uninstall  卸载 115Quick"
    echo "  service    创建并启动 macOS 服务"
    echo "  start      启动服务"
    echo "  stop       停止服务"
    echo "  restart    重启服务"
    echo "  status     显示状态信息"
    echo "  help       显示此帮助信息"
    echo ""
}

# 主逻辑
main() {
    case "${1:-help}" in
        install)
            if [ -f "${INSTALL_DIR}/115quick" ]; then
                info "检测到已安装，执行更新..."
                check_update
            else
                local version
                version=$(get_latest_version)
                install_binary "$version"
                echo ""
                read -p "是否创建 macOS 服务（开机自启）？(Y/n): " -n 1 -r
                echo ""
                if [[ ! $REPLY =~ ^[Nn]$ ]]; then
                    create_service
                    start_service
                else
                    info "跳过服务创建"
                    info "手动启动: 115quick -f ${CONFIG_DIR}/quick.yaml"
                fi
            fi
            ;;
        update)
            check_update
            ;;
        uninstall)
            uninstall
            ;;
        service)
            create_service
            start_service
            ;;
        start)
            if [ ! -f "$PLIST_PATH" ]; then
                error "服务未安装，请先运行: $0 service"
            fi
            launchctl start "$PLIST_NAME"
            success "服务已启动"
            ;;
        stop)
            if [ ! -f "$PLIST_PATH" ]; then
                error "服务未安装"
            fi
            launchctl stop "$PLIST_NAME"
            success "服务已停止"
            ;;
        restart)
            if [ ! -f "$PLIST_PATH" ]; then
                error "服务未安装，请先运行: $0 service"
            fi
            launchctl stop "$PLIST_NAME" 2>/dev/null || true
            launchctl start "$PLIST_NAME"
            success "服务已重启"
            ;;
        status)
            show_status
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "未知命令: $1"
            ;;
    esac
}

main "$@"
