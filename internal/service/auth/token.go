package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"115Quick_server/internal/store"
)

const (
	TokenRefreshURL = "https://passportapi.115.com/open/refreshToken"
	TokenBuffer     = 10 * time.Minute
)

type Manager struct {
	mu           sync.RWMutex
	store        *store.Store
	accessToken  string
	refreshToken string
	expiresAt    time.Time
	lastRefresh  time.Time
	httpClient   *http.Client
	onTokenChange func(accessToken, refreshToken string)
}

type TokenResponse struct {
	State   bool   `json:"state"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	} `json:"data"`
	Error string `json:"error"`
	ErrNo int    `json:"errno"`
}

func NewManager(s *store.Store) *Manager {
	m := &Manager{
		store: s,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	m.loadFromStore()
	return m
}

func (m *Manager) loadFromStore() {
	tc, err := m.store.GetToken()
	if err != nil {
		return
	}

	m.accessToken = tc.AccessToken
	m.refreshToken = tc.RefreshToken
	if tc.ExpiresAt > 0 {
		m.expiresAt = time.Unix(tc.ExpiresAt, 0)
	}
}

func (m *Manager) SetOnTokenChange(fn func(accessToken, refreshToken string)) {
	m.onTokenChange = fn
}

func (m *Manager) SetToken(accessToken, refreshToken string, expiresIn int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accessToken = accessToken
	m.refreshToken = refreshToken
	m.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	m.lastRefresh = time.Now()

	if err := m.store.SaveToken(accessToken, refreshToken, m.expiresAt.Unix()); err != nil {
		return fmt.Errorf("save token to store: %w", err)
	}

	if m.onTokenChange != nil {
		m.onTokenChange(accessToken, refreshToken)
	}

	return nil
}

func (m *Manager) GetAccessToken() (string, error) {
	m.mu.RLock()
	if m.accessToken != "" && time.Now().Add(TokenBuffer).Before(m.expiresAt) {
		token := m.accessToken
		m.mu.RUnlock()
		return token, nil
	}
	m.mu.RUnlock()

	return m.accessToken, m.refresh()
}

func (m *Manager) refresh() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.accessToken != "" && time.Now().Add(TokenBuffer).Before(m.expiresAt) {
		return nil
	}

	if m.refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	return m.doRefresh()
}

func (m *Manager) doRefresh() error {
	data := url.Values{}
	data.Set("refresh_token", m.refreshToken)

	resp, err := m.httpClient.Post(TokenRefreshURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("refresh token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read refresh response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("parse refresh response: %w", err)
	}

	if !tokenResp.State || tokenResp.Code != 0 {
		return fmt.Errorf("refresh failed: code=%d message=%s", tokenResp.Code, tokenResp.Message)
	}

	m.accessToken = tokenResp.Data.AccessToken
	if tokenResp.Data.RefreshToken != "" {
		m.refreshToken = tokenResp.Data.RefreshToken
	}
	m.expiresAt = time.Now().Add(time.Duration(tokenResp.Data.ExpiresIn) * time.Second)
	m.lastRefresh = time.Now()

	if err := m.store.SaveToken(m.accessToken, m.refreshToken, m.expiresAt.Unix()); err != nil {
		return fmt.Errorf("save refreshed token: %w", err)
	}

	return nil
}

func (m *Manager) ForceRefresh() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	return m.doRefresh()
}

func (m *Manager) GetStatus() *TokenStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := &TokenStatus{
		Configured: m.refreshToken != "",
	}

	if m.accessToken != "" {
		status.Valid = time.Now().Before(m.expiresAt)
		status.ExpiresAt = m.expiresAt.Unix()
	}

	if !status.Configured {
		status.Message = "未配置 Token，请通过插件输入 RefreshToken"
	} else if !status.Valid {
		status.Message = "Token 已过期，请刷新或重新输入"
	} else {
		remaining := time.Until(m.expiresAt)
		status.Message = fmt.Sprintf("Token 有效，剩余 %s", remaining.Round(time.Minute))
	}

	return status
}

type TokenStatus struct {
	Configured bool   `json:"configured"`
	Valid      bool   `json:"valid"`
	ExpiresAt  int64  `json:"expires_at"`
	Message    string `json:"message"`
}

func (m *Manager) HasToken() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.refreshToken != ""
}

func (m *Manager) ClearToken() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accessToken = ""
	m.refreshToken = ""
	m.expiresAt = time.Time{}
	m.store.SaveToken("", "", 0)
}
