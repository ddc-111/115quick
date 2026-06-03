package smb

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SMBConfig struct {
	Host       string `json:"host"`
	Share      string `json:"share"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	MountPoint string `json:"mountPoint"`
}

type SMBManager struct {
	config    *SMBConfig
	isMounted bool
}

func NewSMBManager() *SMBManager {
	return &SMBManager{
		isMounted: false,
	}
}

func (m *SMBManager) Mount(cfg *SMBConfig) error {
	if cfg.Host == "" || cfg.Share == "" {
		return fmt.Errorf("SMB host and share are required")
	}

	if cfg.MountPoint == "" {
		return fmt.Errorf("mount point is required")
	}

	if err := os.MkdirAll(cfg.MountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %v", err)
	}

	if m.isMounted {
		if err := m.Unmount(); err != nil {
			return fmt.Errorf("failed to unmount existing share: %v", err)
		}
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		mountArgs := []string{"-t", "cifs"}
		options := fmt.Sprintf("username=%s,password=%s,vers=3.0", cfg.Username, cfg.Password)
		if cfg.Username == "" {
			options = "guest,vers=3.0"
		}
		mountArgs = append(mountArgs, "-o", options)
		uncPath := fmt.Sprintf("//%s/%s", cfg.Host, cfg.Share)
		mountArgs = append(mountArgs, uncPath, cfg.MountPoint)
		cmd = exec.Command("mount", mountArgs...)
	case "windows":
		uncPath := fmt.Sprintf("\\\\%s\\%s", cfg.Host, cfg.Share)
		args := []string{"use", cfg.MountPoint, uncPath}
		if cfg.Username != "" {
			args = append(args, fmt.Sprintf("/user:%s", cfg.Username))
			if cfg.Password != "" {
				args = append(args, cfg.Password)
			}
		}
		cmd = exec.Command("net", args...)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mount failed: %s, error: %v", string(output), err)
	}

	m.config = cfg
	m.isMounted = true
	return nil
}

func (m *SMBManager) Unmount() error {
	if !m.isMounted || m.config == nil {
		return nil
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("umount", m.config.MountPoint)
	case "windows":
		cmd = exec.Command("net", "use", m.config.MountPoint, "/delete", "/y")
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "not mounted") || strings.Contains(string(output), "not found") {
			m.isMounted = false
			return nil
		}
		return fmt.Errorf("unmount failed: %s, error: %v", string(output), err)
	}

	m.isMounted = false
	return nil
}

func (m *SMBManager) IsMounted() bool {
	return m.isMounted
}

func (m *SMBManager) GetConfig() *SMBConfig {
	return m.config
}

func (m *SMBManager) TestConnection(cfg *SMBConfig) error {
	if cfg.Host == "" || cfg.Share == "" {
		return fmt.Errorf("SMB host and share are required")
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		uncPath := fmt.Sprintf("//%s/%s", cfg.Host, cfg.Share)
		args := []string{"-t", "cifs", "-o", "guest,vers=3.0"}
		if cfg.Username != "" {
			args = append(args, "-o", fmt.Sprintf("username=%s,password=%s,vers=3.0", cfg.Username, cfg.Password))
		}
		tempMount, err := os.MkdirTemp("", "smb-test-*")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %v", err)
		}
		defer os.Remove(tempMount)

		args = append(args, uncPath, tempMount)
		cmd := exec.Command("mount", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("connection test failed: %s", string(output))
		}
		exec.Command("umount", tempMount).Run()
		return nil
	case "windows":
		uncPath := fmt.Sprintf("\\\\%s\\%s", cfg.Host, cfg.Share)
		args := []string{"use", uncPath}
		if cfg.Username != "" {
			args = append(args, fmt.Sprintf("/user:%s", cfg.Username))
			if cfg.Password != "" {
				args = append(args, cfg.Password)
			}
		}
		cmd := exec.Command("net", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("connection test failed: %s", string(output))
		}
		exec.Command("net", "use", uncPath, "/delete", "/y").Run()
		return nil
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}