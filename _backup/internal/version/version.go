package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

const (
	GitHubRepo    = "AnJiaHa/115Quick_server"
	CheckURL      = "https://api.github.com/repos/" + GitHubRepo + "/releases/latest"
	UpdatePageURL = "https://github.com/" + GitHubRepo + "/releases/latest"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Body    string `json:"body"`
}

type VersionInfo struct {
	Version      string `json:"version"`
	GitCommit    string `json:"gitCommit"`
	BuildTime    string `json:"buildTime"`
	UpdateURL    string `json:"updateUrl,omitempty"`
	LatestVersion string `json:"latestVersion,omitempty"`
	HasUpdate    bool   `json:"hasUpdate"`
}

func GetVersionInfo() VersionInfo {
	info := VersionInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
	}

	latest, url, err := CheckLatestVersion()
	if err == nil && latest != "" {
		info.LatestVersion = latest
		info.UpdateURL = url
		info.HasUpdate = compareVersions(latest, Version) > 0
	}

	return info
}

func CheckLatestVersion() (string, string, error) {
	client := &http.Client{Timeout: 5 * 1000000000} // 5 seconds
	req, err := http.NewRequest("GET", CheckURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	tag := release.TagName
	if strings.HasPrefix(tag, "v") {
		tag = strings.TrimPrefix(tag, "v")
	}

	url := release.HTMLURL
	if url == "" {
		url = UpdatePageURL
	}

	return tag, url, nil
}

func compareVersions(v1, v2 string) int {
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")

	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(s1) {
			fmt.Sscanf(s1[i], "%d", &n1)
		}
		if i < len(s2) {
			fmt.Sscanf(s2[i], "%d", &n2)
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}
