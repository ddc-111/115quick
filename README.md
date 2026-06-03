# 115Quick

115Quick 是一个基于 115 网盘 API 的离线下载服务，允许用户在自己的服务器上部署，通过 115 网盘进行磁力链接的离线下载，并自动将下载完成的文件保存到指定位置。

## 功能特点

- 🧲 支持磁力链接离线下载
- 🔄 自动轮询检查下载状态
- 📥 自动下载完成文件到指定位置
- 🌐 支持公网和内网访问
- ⚡ 简单易用的 API 接口
- 🖥️ 支持 Windows、Linux、Mac 等多平台
- 🔌 提供 Chrome 插件支持
- 💾 SQLite 本地数据存储
- 🗂️ 支持 SMB 网络存储（NAS、Windows 共享等）
- 🔍 **页面磁力链自动检测** - 浏览任意网页时自动识别磁力链接
- ⬇️ **一键推送下载** - 在磁力链接旁显示悬浮按钮，点击即推送到 115 离线下载
- 🔔 **推送状态提示** - 推送成功/失败时显示 Toast 通知

## 安装使用

### 从 Release 下载

1. 从 [Release](https://github.com/ddc-111/115quick/releases) 页面下载对应平台的二进制文件
2. 解压后修改 `etc/quick.yaml` 配置文件

### 配置

编辑 `etc/quick.yaml`：

```yaml
Name: quick
Host: 0.0.0.0
Port: 8889
DBPath: data/115quick.db
Auth115:
  DownloadPath: data
  AccessToken: "你的AccessToken"
  RefreshToken: "你的RefreshToken"
SMB:
  Enabled: false
  Host: ""
  Share: ""
  Username: ""
  Password: ""
  MountPoint: ""
```

### SMB 网络存储配置

SMB 功能允许将下载文件直接保存到网络共享存储（如 NAS、Windows 共享文件夹等）。

#### 通过 Chrome 插件配置（推荐）

1. 打开 Chrome 插件，进入「设置」页面
2. 找到「SMB 网络存储」配置区域
3. 开启 SMB 开关
4. 填写配置信息：
   - **SMB 服务器地址**: NAS 或共享服务器的 IP 地址或主机名
   - **共享名称**: 共享文件夹名称
   - **用户名**: 访问共享的用户名（匿名访问可留空）
   - **密码**: 访问密码（无密码可留空）
   - **挂载点**: 
     - Linux/Mac: 本地目录路径，如 `/mnt/smb` 或 `/Volumes/share`
     - Windows: 盘符，如 `Z:` 或 `Y:`
5. 点击「测试连接」验证配置是否正确
6. 点击「保存配置」应用设置

#### 通过配置文件配置

```yaml
SMB:
  Enabled: true
  Host: "192.168.1.100"
  Share: "downloads"
  Username: "user"
  Password: "password"
  MountPoint: "/mnt/smb"
```

#### 系统要求

- **Linux**: 需要安装 `cifs-utils` 包 (`sudo apt install cifs-utils`)
- **macOS**: 系统自带 SMB 支持
- **Windows**: 系统自带 SMB 支持，使用盘符作为挂载点

#### 注意事项

1. 确保服务器有权限访问 SMB 共享
2. 挂载点目录会自动创建
3. 如果 SMB 连接失败，下载将使用默认的本地路径
4. 配置保存在数据库中，重启服务后会自动尝试挂载

### 获取 Token

1. 登录 115 网盘开放平台获取 OAuth 应用的 AccessToken 和 RefreshToken
2. 将获取到的 Token 填入配置文件

### 启动服务

```bash
# Windows
115quick.exe

# Linux/Mac
./115quick

# 指定配置文件
./115quick -f etc/quick.yaml
```

### 安装 Chrome 插件

1. 从 Release 下载 Chrome 插件压缩包并解压
2. 打开 Chrome 浏览器，进入 `chrome://extensions/`
3. 开启「开发者模式」
4. 点击「加载已解压的扩展程序」，选择解压后的目录
5. 在插件选项中设置服务器地址（如 `http://localhost:8889`）

### Chrome 插件功能

#### 管理面板

点击浏览器工具栏的 115Quick 图标打开管理面板，包含：

- **仪表盘**: 查看服务器状态、Token 状态、下载进度
- **下载列表**: 查看当前下载任务和云下载任务
- **任务历史**: 查看历史下载记录
- **Token 配置**: 设置 115 网盘 AccessToken 和 RefreshToken
- **设置**: 配置下载模式、SMB 存储、重命名工具

#### 页面磁力链检测

插件会自动在所有网页上检测磁力链接：

1. **自动扫描**: 页面加载时自动扫描所有磁力链接
2. **动态监听**: 页面动态加载内容时自动检测新增的磁力链
3. **悬浮按钮**: 在检测到的磁力链接旁边显示「⬇ 推送下载」按钮
4. **一键推送**: 点击按钮立即将磁力链接推送到 115 离线下载
5. **状态反馈**: 
   - 推送中：按钮显示「推送中...」
   - 成功：按钮变为绿色显示「已推送」，右上角弹出绿色 Toast
   - 失败：按钮显示「重试」，右上角弹出红色 Toast 错误信息

#### 使用示例

在任意包含磁力链接的网页上（如 BT 种子网站、论坛等）：

```
magnet:?xt=urn:btih:abcdef1234567890... [⬇ 推送下载]
```

点击「推送下载」按钮后：
- ✅ 成功：「磁力链接已推送到115下载」
- ❌ 失败：「推送失败: 错误信息」

## API 接口

所有接口前缀：`/v1/Download`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/addDownloadLink` | 提交磁力链接 |
| GET | `/api/getServerInfo` | 获取服务器状态 |
| POST | `/api/setDownloadMode` | 设置下载模式 (0=仅视频, 1=全部) |
| POST | `/api/StartReName` | 触发文件重命名 |
| POST | `/api/setToken` | 设置 115 AccessToken 和 RefreshToken |
| GET | `/api/getTokenStatus` | 获取 Token 状态 |
| GET | `/api/getDownloadProgress` | 获取当前下载进度 |
| GET | `/api/getCloudTasks` | 获取云下载任务列表 |
| GET | `/api/getTaskHistory` | 获取任务历史（分页） |
| POST | `/api/removeDownloadTask` | 删除待下载任务 |
| POST | `/api/refreshTasks` | 手动刷新云任务列表 |
| GET | `/api/getSMBConfig` | 获取 SMB 配置 |
| POST | `/api/setSMBConfig` | 设置 SMB 配置 |
| POST | `/api/testSMBConnection` | 测试 SMB 连接 |

### 示例

```bash
# 添加下载任务
curl -X POST http://localhost:8889/v1/Download/api/addDownloadLink \
  -H "Content-Type: application/json" \
  -d '{"downloadLink": "magnet:?xt=urn:btih:..."}'

# 获取服务器状态
curl http://localhost:8889/v1/Download/api/getServerInfo

# 配置 SMB
curl -X POST http://localhost:8889/v1/Download/api/setSMBConfig \
  -H "Content-Type: application/json" \
  -d '{"enabled": true, "host": "192.168.1.100", "share": "downloads", "username": "user", "password": "pass", "mountPoint": "/mnt/smb"}'

# 测试 SMB 连接
curl -X POST http://localhost:8889/v1/Download/api/testSMBConnection \
  -H "Content-Type: application/json" \
  -d '{"host": "192.168.1.100", "share": "downloads", "username": "user", "password": "pass"}'
```

## 从源码构建

### 环境要求

- Go 1.24+
- Node.js 20+（构建 Chrome 插件）

### 构建服务端

```bash
go build -o 115quick .
```

### 构建 Chrome 插件

```bash
cd plugin/115Quick
npm install
npm run build
# 构建产物在 dist 目录
```

## 项目结构

```
.
├── quick.go                    # 入口文件
├── etc/
│   └── quick.yaml              # 配置文件
├── internal/
│   ├── config/                 # 配置结构
│   ├── handler/                # HTTP 处理器
│   ├── logic/                  # 业务逻辑
│   ├── svc/                    # 核心服务
│   │   ├── servicecontext.go   # 服务上下文
│   │   ├── auth115.go          # 115 API 客户端
│   │   └── store.go            # SQLite 存储
│   ├── types/                  # 类型定义
│   └── utils/
│       └── smb/                # SMB 网络存储工具
├── plugin/
│   └── 115Quick/               # Chrome 插件源码
│       ├── src/
│       │   ├── content/        # Content Script (页面磁力链检测)
│       │   ├── background/     # Background Script (消息处理)
│       │   ├── views/          # 管理面板页面
│       │   ├── api/            # API 接口封装
│       │   └── stores/         # 状态管理
│       └── manifest.json       # 插件配置
└── .github/
    └── workflows/
        └── release.yml         # GitHub Actions 发布流水线
```

## 注意事项

1. 请确保服务器有足够的存储空间
2. 请遵守 115 网盘的使用条款
3. 所有数据存储在本地 SQLite 数据库中
4. 使用 SMB 网络存储时，请确保网络连接稳定
5. Chrome 插件需要在「设置」中配置服务器地址后才能使用页面推送功能

## 许可证

[MIT License](LICENSE)
