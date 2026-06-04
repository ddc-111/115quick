package smb

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hirochachacha/go-smb2"
)

type SMBConfig struct {
	Host     string `json:"host"`
	Share    string `json:"share"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SMBFile struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"isDir"`
	ModTime time.Time `json:"modTime"`
}

type SMBManager struct {
	config    *SMBConfig
	session   *smb2.Session
	share     *smb2.Share
	connected bool
}

func NewSMBManager() *SMBManager {
	return &SMBManager{
		connected: false,
	}
}

func (m *SMBManager) Connect(cfg *SMBConfig) error {
	if cfg.Host == "" || cfg.Share == "" {
		return fmt.Errorf("SMB服务器地址和共享名称不能为空")
	}

	// 断开现有连接
	if m.connected {
		m.Disconnect()
	}

	// 尝试连接
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:445", cfg.Host), 10*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接到SMB服务器 %s:445，请检查：1.服务器IP是否正确 2.SMB服务是否开启 3.防火墙是否允许445端口 (错误: %v)", cfg.Host, err)
	}

	// 创建 SMB dialer
	// Windows 共享无密码时，密码设为空字符串
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     cfg.Username,
			Password: cfg.Password,
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		conn.Close()
		return fmt.Errorf("SMB会话建立失败: %v。Windows共享请确认：1.用户名是否正确 2.共享是否开启密码保护", err)
	}

	share, err := session.Mount(cfg.Share)
	if err != nil {
		session.Logoff()
		return fmt.Errorf("挂载SMB共享 '%s' 失败: %v。请检查共享名称是否正确（如 C、D、Downloads）", cfg.Share, err)
	}

	m.config = cfg
	m.session = session
	m.share = share
	m.connected = true

	return nil
}

func (m *SMBManager) TestConnection(cfg *SMBConfig) error {
	if cfg.Host == "" || cfg.Share == "" {
		return fmt.Errorf("SMB服务器地址和共享名称不能为空")
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:445", cfg.Host), 10*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接到SMB服务器 %s:445: %v", cfg.Host, err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     cfg.Username,
			Password: cfg.Password,
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		return fmt.Errorf("SMB会话建立失败: %v", err)
	}
	defer session.Logoff()

	share, err := session.Mount(cfg.Share)
	if err != nil {
		return fmt.Errorf("挂载共享 '%s' 失败: %v", cfg.Share, err)
	}
	defer share.Umount()

	return nil
}

func (m *SMBManager) Disconnect() error {
	if m.share != nil {
		m.share.Umount()
		m.share = nil
	}
	if m.session != nil {
		m.session.Logoff()
		m.session = nil
	}
	m.connected = false
	return nil
}

func (m *SMBManager) IsConnected() bool {
	return m.connected
}

func (m *SMBManager) GetConfig() *SMBConfig {
	return m.config
}

func (m *SMBManager) ListFiles(path string) ([]SMBFile, error) {
	if !m.connected || m.share == nil {
		return nil, fmt.Errorf("not connected to SMB server")
	}

	if path == "" {
		path = "."
	}

	entries, err := m.share.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	files := make([]SMBFile, 0, len(entries))
	for _, entry := range entries {
		files = append(files, SMBFile{
			Name:    entry.Name(),
			Size:    entry.Size(),
			IsDir:   entry.IsDir(),
			ModTime: entry.ModTime(),
		})
	}

	return files, nil
}

func (m *SMBManager) DownloadFile(remotePath, localPath string) error {
	if !m.connected || m.share == nil {
		return fmt.Errorf("not connected to SMB server")
	}

	// 确保本地目录存在
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %v", err)
	}

	// 打开远程文件
	srcFile, err := m.share.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer srcFile.Close()

	// 创建本地文件
	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	return nil
}

func (m *SMBManager) GetFileInfo(path string) (*SMBFile, error) {
	if !m.connected || m.share == nil {
		return nil, fmt.Errorf("not connected to SMB server")
	}

	info, err := m.share.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	return &SMBFile{
		Name:    info.Name(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
	}, nil
}
