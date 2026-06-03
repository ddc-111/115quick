package svc

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS file_infos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		info_hash TEXT NOT NULL DEFAULT '',
		url TEXT NOT NULL,
		add_time DATETIME NOT NULL
	)`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS token_config (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		access_token TEXT NOT NULL,
		refresh_token TEXT NOT NULL,
		updated_at DATETIME NOT NULL
	)`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS task_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		info_hash TEXT DEFAULT '',
		name TEXT DEFAULT '',
		size INTEGER DEFAULT 0,
		status INTEGER NOT NULL DEFAULT 0,
		progress REAL DEFAULT 0,
		error_msg TEXT DEFAULT '',
		add_time DATETIME NOT NULL,
		complete_time DATETIME
	)`); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS smb_config (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		enabled INTEGER NOT NULL DEFAULT 0,
		host TEXT DEFAULT '',
		share TEXT DEFAULT '',
		username TEXT DEFAULT '',
		password TEXT DEFAULT '',
		mount_point TEXT DEFAULT '',
		updated_at DATETIME NOT NULL
	)`); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) AddFileInfo(infoHash, url string) error {
	_, err := s.db.Exec("INSERT INTO file_infos (info_hash, url, add_time) VALUES (?, ?, ?)",
		infoHash, url, time.Now())
	return err
}

func (s *Store) RemoveFileInfo(url string) error {
	_, err := s.db.Exec("DELETE FROM file_infos WHERE url = ?", url)
	return err
}

func (s *Store) GetAllFileInfos() ([]FileInfo, error) {
	rows, err := s.db.Query("SELECT info_hash, url, add_time FROM file_infos ORDER BY add_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var infos []FileInfo
	for rows.Next() {
		var fi FileInfo
		if err := rows.Scan(&fi.InfoHash, &fi.URL, &fi.AddTime); err != nil {
			return nil, err
		}
		infos = append(infos, fi)
	}
	return infos, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// Token 相关方法

func (s *Store) SaveToken(accessToken, refreshToken string) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO token_config (id, access_token, refresh_token, updated_at) 
		VALUES (1, ?, ?, ?)`, accessToken, refreshToken, time.Now())
	return err
}

func (s *Store) GetToken() (accessToken, refreshToken string, err error) {
	row := s.db.QueryRow("SELECT access_token, refresh_token FROM token_config WHERE id = 1")
	err = row.Scan(&accessToken, &refreshToken)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	return
}

// 任务历史相关方法

func (s *Store) AddTaskHistory(url, infoHash, name string, size int64, status int) error {
	_, err := s.db.Exec(`INSERT INTO task_history (url, info_hash, name, size, status, add_time) 
		VALUES (?, ?, ?, ?, ?, ?)`, url, infoHash, name, size, status, time.Now())
	return err
}

func (s *Store) UpdateTaskHistoryStatus(url string, status int, progress float64, errorMsg string) error {
	_, err := s.db.Exec(`UPDATE task_history SET status = ?, progress = ?, error_msg = ? WHERE url = ?`,
		status, progress, errorMsg, url)
	return err
}

func (s *Store) UpdateTaskHistoryComplete(url string, status int) error {
	_, err := s.db.Exec(`UPDATE task_history SET status = ?, complete_time = ? WHERE url = ?`,
		status, time.Now(), url)
	return err
}

func (s *Store) GetTaskHistory(page, pageSize int) (items []struct {
	ID           int64
	URL          string
	Name         string
	Size         int64
	Status       int
	Progress     float64
	ErrorMsg     string
	AddTime      time.Time
	CompleteTime *time.Time
}, total int, err error) {
	// 获取总数
	err = s.db.QueryRow("SELECT COUNT(*) FROM task_history").Scan(&total)
	if err != nil {
		return
	}

	// 分页查询
	offset := (page - 1) * pageSize
	rows, err := s.db.Query(`SELECT id, url, name, size, status, progress, error_msg, add_time, complete_time 
		FROM task_history ORDER BY add_time DESC LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item struct {
			ID           int64
			URL          string
			Name         string
			Size         int64
			Status       int
			Progress     float64
			ErrorMsg     string
			AddTime      time.Time
			CompleteTime *time.Time
		}
		if err = rows.Scan(&item.ID, &item.URL, &item.Name, &item.Size, &item.Status,
			&item.Progress, &item.ErrorMsg, &item.AddTime, &item.CompleteTime); err != nil {
			return
		}
		items = append(items, item)
	}
	return
}

func (s *Store) RemoveTaskHistory(url string) error {
	_, err := s.db.Exec("DELETE FROM task_history WHERE url = ?", url)
	return err
}

func (s *Store) ClearTaskHistory() error {
	_, err := s.db.Exec("DELETE FROM task_history")
	return err
}

// SMB 配置相关方法

type SMBConfigData struct {
	Enabled    bool
	Host       string
	Share      string
	Username   string
	Password   string
	MountPoint string
}

func (s *Store) SaveSMBConfig(cfg SMBConfigData) error {
	enabled := 0
	if cfg.Enabled {
		enabled = 1
	}
	_, err := s.db.Exec(`INSERT OR REPLACE INTO smb_config (id, enabled, host, share, username, password, mount_point, updated_at) 
		VALUES (1, ?, ?, ?, ?, ?, ?, ?)`, enabled, cfg.Host, cfg.Share, cfg.Username, cfg.Password, cfg.MountPoint, time.Now())
	return err
}

func (s *Store) GetSMBConfig() (*SMBConfigData, error) {
	row := s.db.QueryRow("SELECT enabled, host, share, username, password, mount_point FROM smb_config WHERE id = 1")
	var cfg SMBConfigData
	var enabled int
	err := row.Scan(&enabled, &cfg.Host, &cfg.Share, &cfg.Username, &cfg.Password, &cfg.MountPoint)
	if err == sql.ErrNoRows {
		return &SMBConfigData{}, nil
	}
	if err != nil {
		return nil, err
	}
	cfg.Enabled = enabled == 1
	return &cfg, nil
}
