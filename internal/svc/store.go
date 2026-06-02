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
