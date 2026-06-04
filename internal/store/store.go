package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS token_config (
			id INTEGER PRIMARY KEY DEFAULT 1,
			access_token TEXT NOT NULL DEFAULT '',
			refresh_token TEXT NOT NULL DEFAULT '',
			expires_at INTEGER NOT NULL DEFAULT 0,
			updated_at INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS task_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL DEFAULT '',
			info_hash TEXT NOT NULL DEFAULT '',
			name TEXT NOT NULL DEFAULT '',
			size INTEGER NOT NULL DEFAULT 0,
			status INTEGER NOT NULL DEFAULT 0,
			progress REAL NOT NULL DEFAULT 0,
			error_msg TEXT NOT NULL DEFAULT '',
			add_time INTEGER NOT NULL DEFAULT 0,
			complete_time INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS file_infos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			info_hash TEXT NOT NULL DEFAULT '',
			url TEXT NOT NULL DEFAULT '',
			add_time INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL DEFAULT ''
		)`,
		`INSERT OR IGNORE INTO token_config (id) VALUES (1)`,
	}

	for _, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("execute migration: %w", err)
		}
	}

	return nil
}

// Token operations
type TokenConfig struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
	UpdatedAt    int64
}

func (s *Store) SaveToken(accessToken, refreshToken string, expiresAt int64) error {
	_, err := s.db.Exec(
		`UPDATE token_config SET access_token=?, refresh_token=?, expires_at=?, updated_at=? WHERE id=1`,
		accessToken, refreshToken, expiresAt, time.Now().Unix(),
	)
	return err
}

func (s *Store) GetToken() (*TokenConfig, error) {
	tc := &TokenConfig{}
	err := s.db.QueryRow(
		`SELECT access_token, refresh_token, expires_at, updated_at FROM token_config WHERE id=1`,
	).Scan(&tc.AccessToken, &tc.RefreshToken, &tc.ExpiresAt, &tc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

// Task history operations
type TaskHistory struct {
	ID           int64
	URL          string
	InfoHash     string
	Name         string
	Size         int64
	Status       int
	Progress     float64
	ErrorMsg     string
	AddTime      int64
	CompleteTime int64
}

func (s *Store) AddTaskHistory(url, infoHash, name string, size int64) error {
	_, err := s.db.Exec(
		`INSERT INTO task_history (url, info_hash, name, size, status, add_time) VALUES (?, ?, ?, ?, 0, ?)`,
		url, infoHash, name, size, time.Now().Unix(),
	)
	return err
}

func (s *Store) UpdateTaskHistoryStatus(infoHash string, status int, errorMsg string) error {
	_, err := s.db.Exec(
		`UPDATE task_history SET status=?, error_msg=? WHERE info_hash=?`,
		status, errorMsg, infoHash,
	)
	return err
}

func (s *Store) UpdateTaskHistoryComplete(infoHash string, status int) error {
	_, err := s.db.Exec(
		`UPDATE task_history SET status=?, complete_time=? WHERE info_hash=?`,
		status, time.Now().Unix(), infoHash,
	)
	return err
}

func (s *Store) GetTaskHistory(page, pageSize int) ([]TaskHistory, int, error) {
	var total int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM task_history`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.Query(
		`SELECT id, url, info_hash, name, size, status, progress, error_msg, add_time, complete_time 
		 FROM task_history ORDER BY id DESC LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []TaskHistory
	for rows.Next() {
		var t TaskHistory
		if err := rows.Scan(&t.ID, &t.URL, &t.InfoHash, &t.Name, &t.Size, &t.Status, &t.Progress, &t.ErrorMsg, &t.AddTime, &t.CompleteTime); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}

	return tasks, total, nil
}

func (s *Store) RemoveTaskHistory(id int64) error {
	_, err := s.db.Exec(`DELETE FROM task_history WHERE id=?`, id)
	return err
}

func (s *Store) ClearTaskHistory() error {
	_, err := s.db.Exec(`DELETE FROM task_history`)
	return err
}

func (s *Store) ClearCompletedTasks() error {
	_, err := s.db.Exec(`DELETE FROM task_history WHERE status=2`)
	return err
}

// File info operations
type FileInfo struct {
	ID       int64
	InfoHash string
	URL      string
	AddTime  int64
}

func (s *Store) AddFileInfo(infoHash, url string) error {
	_, err := s.db.Exec(
		`INSERT INTO file_infos (info_hash, url, add_time) VALUES (?, ?, ?)`,
		infoHash, url, time.Now().Unix(),
	)
	return err
}

func (s *Store) RemoveFileInfo(infoHash string) error {
	_, err := s.db.Exec(`DELETE FROM file_infos WHERE info_hash=?`, infoHash)
	return err
}

func (s *Store) GetAllFileInfos() ([]FileInfo, error) {
	rows, err := s.db.Query(`SELECT id, info_hash, url, add_time FROM file_infos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var infos []FileInfo
	for rows.Next() {
		var fi FileInfo
		if err := rows.Scan(&fi.ID, &fi.InfoHash, &fi.URL, &fi.AddTime); err != nil {
			return nil, err
		}
		infos = append(infos, fi)
	}
	return infos, nil
}

// Config operations
func (s *Store) SetConfig(key, value string) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)`,
		key, value,
	)
	return err
}

func (s *Store) GetConfig(key string) (string, error) {
	var value string
	err := s.db.QueryRow(`SELECT value FROM config WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
