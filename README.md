# 115Quick

115Quick 是一个基于 115 网盘 API 的离线下载服务，允许用户在自己的服务器上部署，通过 115 网盘进行磁力链接的离线下载，并自动将下载完成的文件保存到指定位置。

## 功能特点

- 🔐 支持 115 网盘账号认证
- 🧲 支持磁力链接离线下载
- 🔄 自动轮询检查下载状态
- 📥 自动下载完成文件到指定位置
- 🌐 支持公网和内网访问
- ⚡ 简单易用的 API 接口
- 🖥️ 支持 Windows、Linux、Mac、NAS 等多平台
- 🔌 提供 Chrome 插件支持

## 安装使用步骤

### 环境要求

- 115 网盘账号
- 支持 Windows、Linux、Mac、NAS 等操作系统

### 快速开始

1. 从 [Release](https://github.com/ddc-111/115quick/releases/tag/0.1) 页面下载对应平台的二进制文件 在main分钟获取etc配置文件夹
2. 解压文件后，在 `etc` 目录下修改 `config.yaml` 配置文件：
   - 设置下载文件保存路径
   - 配置服务监听端口
3. 添加权限并启动服务：
   - Windows: 双击运行 `115Quick.exe`
   - Linux/Mac: 执行 `./115Quick`
4. 授权配置：
   - 启动后查看控制台输出二维码图片
   - 或在服务同级目录下查看生成的二维码图片
   - 使用手机扫描二维码完成 115 网盘授权
   - 授权数据将安全存储在本地，无需担心隐私问题
5. 安装 Chrome 插件：
   - 下载并解压 Chrome 插件
   - 打开 Chrome 浏览器，进入扩展程序页面
   - 开启开发者模式
   - 加载已解压的扩展程序
6. 配置 Chrome 插件：
   - 在插件选项中设置服务器 IP 和端口
   - 完成配置后即可使用

### 使用说明

1. 复制磁力链接后，插件会自动将链接发送到 115Quick 服务
2. 服务会自动使用 115 网盘进行离线下载
3. 下载完成后，文件会自动保存到配置的下载路径

## 配置说明

在 `config.yaml` 文件中可以配置以下参数：

- 下载文件保存路径
- 服务监听端口

## 注意事项

1. 请确保服务器有足够的存储空间
2. 定期检查日志文件，确保服务正常运行
3. 请遵守 115 网盘的使用条款
4. 本服务完全本地化运行，所有数据存储在本地，无需担心隐私问题

## 密钥检测

- 本服务提供 24 小时免费试用
- 试用期间后可以进qq群966107802获取免费密钥
- ![qrcode_1744876569463](https://github.com/user-attachments/assets/327eab27-a66e-47e8-b0cd-78761e35fa60)


我的爱发电主页 https://afdian.com/a/quick115
![afdian- 未认证 115quick](https://github.com/user-attachments/assets/ac223324-63f7-4d86-83ef-b4d5c502b1ef)
