package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DBPath string
	Auth115 struct {
		DownloadPath string
		AccessToken  string
		RefreshToken string
	}
	SMB struct {
		Enabled  bool   `json:",default=false"`
		Host     string `json:",optional"`
		Share    string `json:",optional"`
		Username string `json:",optional"`
		Password string `json:",optional"`
		MountPoint string `json:",optional"`
	}
}
